<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\Website;
use App\Services\Agent\AgentClientInterface;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Response as HttpResponse;
use Illuminate\Support\Facades\File;
use Illuminate\Support\Str;
use Inertia\Inertia;
use Inertia\Response;
use Symfony\Component\HttpFoundation\StreamedResponse;
use ZipArchive;

class FileManagerController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

    /**
     * Validate that base_path is a legitimate, allowed path inside /var/www.
     */
    protected function sanitizeAndValidateBasePath(string $basePath, ?Website $website = null): string
    {
        if (str_contains($basePath, "\0")) {
            abort(400, 'Invalid path: null byte detected.');
        }

        $normalized = str_replace('\\', '/', $basePath);
        $parts = array_filter(explode('/', $normalized), strlen(...));
        $resolvedParts = [];
        foreach ($parts as $part) {
            if ($part === '.') continue;
            if ($part === '..') {
                array_pop($resolvedParts);
            } else {
                $resolvedParts[] = $part;
            }
        }
        $clean = '/' . implode('/', $resolvedParts);

        $allowedPrefixes = ['/var/www', '/private/var/www', '/tmp', '/private/tmp'];
        $isAllowed = false;
        foreach ($allowedPrefixes as $prefix) {
            if ($clean === $prefix || str_starts_with($clean, $prefix . '/')) {
                $isAllowed = true;
                break;
            }
        }

        if (!$isAllowed) {
            return $website && $website->document_root ? dirname($website->document_root) : '/var/www';
        }

        return $clean;
    }

    public function index(Request $request): Response
    {
        $websites = Website::select('id', 'domain', 'document_root', 'system_user')
            ->orderBy('domain')
            ->get();

        $selectedWebsiteId = $request->query('website_id', $websites->first()?->id);
        $selectedWebsite = $websites->firstWhere('id', $selectedWebsiteId) ?? $websites->first();

        $basePath = '/var/www';
        if ($selectedWebsite && $selectedWebsite->document_root) {
            $basePath = dirname($selectedWebsite->document_root);
        }
        $basePath = $this->sanitizeAndValidateBasePath($basePath, $selectedWebsite);

        $relativePath = (string) $request->query('path', '');
        $showHidden = filter_var($request->query('show_hidden', false), FILTER_VALIDATE_BOOLEAN);

        $files = [];
        $diskUsage = [];
        try {
            $files = $this->agentClient->browseFiles($basePath, $relativePath, $showHidden);
            $diskUsage = $this->agentClient->getDiskUsage($basePath);
        } catch (Exception $e) {
            // Graceful fallback on initial load error
        }

        return Inertia::render('Files/Index', [
            'websites' => $websites,
            'selectedWebsite' => $selectedWebsite,
            'currentPath' => $relativePath,
            'basePath' => $basePath,
            'files' => $files,
            'diskUsage' => $diskUsage,
            'showHidden' => $showHidden,
        ]);
    }

    public function browse(Request $request): JsonResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['nullable', 'string'],
            'show_hidden' => ['nullable', 'boolean'],
            'website_id' => ['nullable'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);
        $relativePath = $validated['relative_path'] ?? '';
        $showHidden = (bool) ($validated['show_hidden'] ?? false);

        try {
            $files = $this->agentClient->browseFiles($basePath, $relativePath, $showHidden);
            $diskUsage = $this->agentClient->getDiskUsage($basePath);

            return response()->json([
                'success' => true,
                'files' => $files,
                'disk_usage' => $diskUsage,
                'base_path' => $basePath,
                'current_path' => $relativePath,
            ]);
        } catch (Exception $e) {
            return response()->json([
                'success' => false,
                'error' => $e->getMessage(),
            ], 400);
        }
    }

    public function read(Request $request): JsonResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $content = $this->agentClient->readFile($basePath, $validated['relative_path']);
            return response()->json([
                'success' => true,
                'content' => $content,
            ]);
        } catch (Exception $e) {
            return response()->json([
                'success' => false,
                'error' => $e->getMessage(),
            ], 400);
        }
    }

    public function write(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
            'content' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $this->agentClient->writeFile($basePath, $validated['relative_path'], $validated['content']);

            $this->logActivity($request, 'file.write', $validated['relative_path'], 'success');

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'File saved successfully.']);
            }

            $url = route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ]);
            return redirect($url)->with('success', 'File saved successfully.');
        } catch (Exception $e) {
            $this->logActivity($request, 'file.write', $validated['relative_path'], 'failure', ['error' => $e->getMessage()]);
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to write file: ' . $e->getMessage());
        }
    }

    public function createFile(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $this->agentClient->createFile($basePath, $validated['relative_path']);
            $this->logActivity($request, 'file.create', $validated['relative_path'], 'success');

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'File created successfully.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'File created successfully.');
        } catch (Exception $e) {
            $this->logActivity($request, 'file.create', $validated['relative_path'], 'failure', ['error' => $e->getMessage()]);
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to create file: ' . $e->getMessage());
        }
    }

    public function createFolder(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $this->agentClient->createDirectory($basePath, $validated['relative_path']);
            $this->logActivity($request, 'file.create_folder', $validated['relative_path'], 'success');

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Folder created successfully.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Folder created successfully.');
        } catch (Exception $e) {
            $this->logActivity($request, 'file.create_folder', $validated['relative_path'], 'failure', ['error' => $e->getMessage()]);
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to create folder: ' . $e->getMessage());
        }
    }

    public function delete(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['nullable', 'string'],
            'paths' => ['nullable', 'array'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);
        $pathsToDelete = $validated['paths'] ?? ($validated['relative_path'] ? [$validated['relative_path']] : []);

        try {
            foreach ($pathsToDelete as $path) {
                if (!empty($path)) {
                    $this->agentClient->deleteFile($basePath, $path);
                    $this->logActivity($request, 'file.delete', $path, 'success');
                }
            }

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Selected item(s) deleted successfully.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Item(s) deleted successfully.');
        } catch (Exception $e) {
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to delete: ' . $e->getMessage());
        }
    }

    public function rename(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'old_path' => ['required', 'string'],
            'new_path' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $this->agentClient->renameFile($basePath, $validated['old_path'], $validated['new_path']);
            $this->logActivity($request, 'file.rename', "{$validated['old_path']} -> {$validated['new_path']}", 'success');

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Renamed successfully.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Renamed successfully.');
        } catch (Exception $e) {
            $this->logActivity($request, 'file.rename', $validated['old_path'], 'failure', ['error' => $e->getMessage()]);
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to rename: ' . $e->getMessage());
        }
    }

    public function copy(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'src_path' => ['nullable', 'string'],
            'sources' => ['nullable', 'array'],
            'dest_path' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);
        $sources = $validated['sources'] ?? ($validated['src_path'] ? [$validated['src_path']] : []);

        try {
            foreach ($sources as $src) {
                $dest = rtrim($validated['dest_path'], '/') . '/' . basename($src);
                $this->agentClient->copyFile($basePath, $src, $dest);
                $this->logActivity($request, 'file.copy', "{$src} -> {$dest}", 'success');
            }

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Copied successfully.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Copied successfully.');
        } catch (Exception $e) {
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to copy: ' . $e->getMessage());
        }
    }

    public function move(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'src_path' => ['nullable', 'string'],
            'sources' => ['nullable', 'array'],
            'dest_path' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);
        $sources = $validated['sources'] ?? ($validated['src_path'] ? [$validated['src_path']] : []);

        try {
            foreach ($sources as $src) {
                $dest = rtrim($validated['dest_path'], '/') . '/' . basename($src);
                $this->agentClient->moveFile($basePath, $src, $dest);
                $this->logActivity($request, 'file.move', "{$src} -> {$dest}", 'success');
            }

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Moved successfully.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Moved successfully.');
        } catch (Exception $e) {
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to move: ' . $e->getMessage());
        }
    }

    public function chmod(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
            'mode' => ['required', 'string'],
            'recursive' => ['nullable', 'boolean'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $this->agentClient->chmodFile(
                $basePath,
                $validated['relative_path'],
                $validated['mode'],
                (bool) ($validated['recursive'] ?? false)
            );
            $this->logActivity($request, 'file.chmod', "{$validated['relative_path']} ({$validated['mode']})", 'success');

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Permissions updated.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Permissions updated.');
        } catch (Exception $e) {
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to update permissions: ' . $e->getMessage());
        }
    }

    public function chown(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
            'uid' => ['required', 'integer'],
            'gid' => ['required', 'integer'],
            'recursive' => ['nullable', 'boolean'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $this->agentClient->chownFile(
                $basePath,
                $validated['relative_path'],
                $validated['uid'],
                $validated['gid'],
                (bool) ($validated['recursive'] ?? false)
            );
            $this->logActivity($request, 'file.chown', $validated['relative_path'], 'success');

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Ownership updated.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Ownership updated.');
        } catch (Exception $e) {
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to update ownership: ' . $e->getMessage());
        }
    }

    public function stat(Request $request): JsonResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $stat = $this->agentClient->statFile($basePath, $validated['relative_path']);
            return response()->json(['success' => true, 'data' => $stat]);
        } catch (Exception $e) {
            return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
        }
    }

    public function compress(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'sources' => ['required', 'array'],
            'dest_path' => ['required', 'string'],
            'format' => ['nullable', 'string', 'in:zip,tar.gz'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);
        $format = $validated['format'] ?? 'zip';

        try {
            $this->agentClient->compressFiles($basePath, $validated['sources'], $validated['dest_path'], $format);
            $this->logActivity($request, 'file.compress', "{$validated['dest_path']}", 'success');

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Archive created successfully.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Archive created successfully.');
        } catch (Exception $e) {
            $this->logActivity($request, 'file.compress', $validated['dest_path'], 'failure', ['error' => $e->getMessage()]);
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to create archive: ' . $e->getMessage());
        }
    }

    public function extract(Request $request): JsonResponse|RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'archive_path' => ['required', 'string'],
            'dest_path' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $this->agentClient->extractArchive($basePath, $validated['archive_path'], $validated['dest_path']);
            $this->logActivity($request, 'file.extract', "{$validated['archive_path']} -> {$validated['dest_path']}", 'success');

            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => true, 'message' => 'Archive extracted successfully.']);
            }

            return redirect()->route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ])->with('success', 'Archive extracted successfully.');
        } catch (Exception $e) {
            $this->logActivity($request, 'file.extract', $validated['archive_path'], 'failure', ['error' => $e->getMessage()]);
            if ($this->isAjaxOrJson($request)) {
                return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
            }
            return back()->with('error', 'Failed to extract archive: ' . $e->getMessage());
        }
    }

    public function search(Request $request): JsonResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'query' => ['required', 'string', 'min:1'],
            'max_results' => ['nullable', 'integer', 'max:500'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);

        try {
            $results = $this->agentClient->searchFiles(
                $basePath,
                $validated['query'],
                $validated['max_results'] ?? 100
            );
            return response()->json(['success' => true, 'data' => $results]);
        } catch (Exception $e) {
            return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
        }
    }

    public function upload(Request $request): JsonResponse
    {
        $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['nullable', 'string'],
            'files' => ['required', 'array'],
            'files.*' => ['file', 'max:102400'], // 100MB max per file
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($request->input('base_path'));
        $targetDir = $request->input('relative_path', '');

        $uploaded = [];
        $errors = [];

        foreach ($request->file('files', []) as $uploadedFile) {
            $filename = $uploadedFile->getClientOriginalName();
            $relPath = $targetDir ? (rtrim($targetDir, '/') . '/' . $filename) : $filename;

            try {
                $rawContent = file_get_contents($uploadedFile->getRealPath());
                $base64Content = base64_encode($rawContent);
                $this->agentClient->writeFileBase64($basePath, $relPath, $base64Content);
                $this->logActivity($request, 'file.upload', $relPath, 'success');
                $uploaded[] = $filename;
            } catch (Exception $e) {
                $errors[] = "{$filename}: " . $e->getMessage();
                $this->logActivity($request, 'file.upload', $relPath, 'failure', ['error' => $e->getMessage()]);
            }
        }

        if (count($uploaded) > 0) {
            return response()->json([
                'success' => true,
                'uploaded' => $uploaded,
                'errors' => $errors,
                'message' => count($uploaded) . ' file(s) uploaded successfully.',
            ]);
        }

        return response()->json([
            'success' => false,
            'errors' => $errors,
            'message' => 'Failed to upload files: ' . implode(', ', $errors),
        ], 400);
    }

    public function download(Request $request): StreamedResponse|HttpResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'paths' => ['nullable'],
            'relative_path' => ['nullable', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);
        $rawPaths = $validated['paths'] ?? ($validated['relative_path'] ? [$validated['relative_path']] : []);
        $paths = is_array($rawPaths) ? $rawPaths : ((is_string($rawPaths) && json_validate($rawPaths)) ? json_decode($rawPaths, true) : [$rawPaths]);

        if (empty($paths)) {
            abort(400, 'No files or folders selected for download.');
        }

        $this->logActivity($request, 'file.download', implode(', ', array_slice($paths, 0, 5)), 'success');

        // Single file direct stream
        if (count($paths) === 1 && !is_dir("{$basePath}/{$paths[0]}")) {
            $filePath = "{$basePath}/{$paths[0]}";
            $filename = basename($paths[0]);

            if (file_exists($filePath) && is_readable($filePath)) {
                return response()->download($filePath, $filename);
            }

            // Fallback to agent read base64
            try {
                $base64 = $this->agentClient->readFileBase64($basePath, $paths[0]);
                if (!empty($base64)) {
                    $raw = base64_decode($base64);
                    return response($raw, 200, [
                        'Content-Type' => 'application/octet-stream',
                        'Content-Disposition' => "attachment; filename=\"{$filename}\"",
                        'Content-Length' => strlen($raw),
                    ]);
                }
            } catch (Exception) {
                // fallback to readFile
            }

            $content = $this->agentClient->readFile($basePath, $paths[0]);
            return response($content, 200, [
                'Content-Type' => 'application/octet-stream',
                'Content-Disposition' => "attachment; filename=\"{$filename}\"",
            ]);
        }

        // Multiple files or folder: Dynamic zip archive stream via Agent or PHP
        $archiveName = (count($paths) === 1 ? basename($paths[0]) : 'archive') . '-' . date('Ymd-His') . '.zip';
        $tempArchive = sys_get_temp_dir() . '/' . Str::uuid() . '.zip';

        try {
            $this->agentClient->compressFiles($basePath, $paths, $tempArchive, 'zip');
            if (file_exists($tempArchive)) {
                return response()->download($tempArchive, $archiveName)->deleteFileAfterSend(true);
            }
        } catch (Exception) {
            // fallback
        }

        return response()->streamDownload(function () use ($basePath, $paths) {
            $zip = new ZipArchive();
            $tempFile = tempnam(sys_get_temp_dir(), 'kp_zip_');
            $zip->open($tempFile, ZipArchive::CREATE | ZipArchive::OVERWRITE);

            foreach ($paths as $relPath) {
                $fullPath = "{$basePath}/{$relPath}";
                if (is_dir($fullPath) && is_readable($fullPath)) {
                    $files = new \RecursiveIteratorIterator(
                        new \RecursiveDirectoryIterator($fullPath, \RecursiveDirectoryIterator::SKIP_DOTS),
                        \RecursiveIteratorIterator::LEAVES_ONLY
                    );
                    foreach ($files as $file) {
                        if (!$file->isDir()) {
                            $filePath = $file->getRealPath();
                            $zipRelPath = substr($filePath, strlen(dirname($fullPath)) + 1);
                            $zip->addFile($filePath, $zipRelPath);
                        }
                    }
                } elseif (file_exists($fullPath) && is_readable($fullPath)) {
                    $zip->addFile($fullPath, basename($fullPath));
                } else {
                    try {
                        $base64 = $this->agentClient->readFileBase64($basePath, $relPath);
                        if (!empty($base64)) {
                            $zip->addFromString(basename($relPath), base64_decode($base64));
                        }
                    } catch (Exception) {
                        // Skip if unreadable
                    }
                }
            }

            $zip->close();
            readfile($tempFile);
            @unlink($tempFile);
        }, $archiveName, [
            'Content-Type' => 'application/zip',
        ]);
    }

    public function preview(Request $request): HttpResponse|JsonResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
        ]);

        $basePath = $this->sanitizeAndValidateBasePath($validated['base_path']);
        $fullPath = "{$basePath}/{$validated['relative_path']}";
        $filename = basename($validated['relative_path']);

        if (file_exists($fullPath) && is_readable($fullPath)) {
            $mime = mime_content_type($fullPath) ?: 'application/octet-stream';
            return response(file_get_contents($fullPath), 200, [
                'Content-Type' => $mime,
                'Content-Disposition' => 'inline; filename="' . $filename . '"',
            ]);
        }

        try {
            $base64 = $this->agentClient->readFileBase64($basePath, $validated['relative_path']);
            if (!empty($base64)) {
                $rawBytes = base64_decode($base64);
                $finfo = new \finfo(FILEINFO_MIME_TYPE);
                $mime = $finfo->buffer($rawBytes) ?: 'application/octet-stream';
                return response($rawBytes, 200, [
                    'Content-Type' => $mime,
                    'Content-Disposition' => 'inline; filename="' . $filename . '"',
                ]);
            }
        } catch (Exception) {
            // fallback
        }

        try {
            $content = $this->agentClient->readFile($basePath, $validated['relative_path']);
            return response($content, 200, [
                'Content-Type' => 'text/plain',
                'Content-Disposition' => 'inline',
            ]);
        } catch (Exception $e) {
            return response()->json(['success' => false, 'error' => $e->getMessage()], 400);
        }
    }

    protected function isAjaxOrJson(Request $request): bool
    {
        return $request->expectsJson()
            || $request->wantsJson()
            || $request->isJson()
            || $request->ajax()
            || str_contains((string) $request->header('Accept'), 'application/json')
            || str_contains((string) $request->header('Content-Type'), 'application/json');
    }

    protected function logActivity(Request $request, string $action, string $resourceId, string $status = 'success', array $payload = []): void
    {
        try {
            ActivityLog::create([
                'user_id' => $request->user()?->id ?? 1,
                'user_email' => $request->user()?->email ?? 'admin@kodepreneur.com',
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => $action,
                'resource_type' => 'file',
                'resource_id' => Str::limit($resourceId, 255),
                'status' => $status,
                'payload_summary' => !empty($payload) ? $payload : null,
            ]);
        } catch (Exception) {
            // Ignore logging failures to prevent blocking core operations
        }
    }
}
