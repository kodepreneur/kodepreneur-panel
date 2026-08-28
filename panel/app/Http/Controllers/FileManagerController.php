<?php

namespace App\Http\Controllers;

use App\Models\ActivityLog;
use App\Models\Website;
use App\Services\Agent\AgentClientInterface;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\RedirectResponse;
use Illuminate\Http\Request;
use Inertia\Inertia;
use Inertia\Response;

class FileManagerController extends Controller
{
    public function __construct(
        protected AgentClientInterface $agentClient
    ) {}

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

        $relativePath = $request->query('path', '');

        $files = [];
        try {
            $files = $this->agentClient->browseFiles($basePath, $relativePath);
        } catch (Exception $e) {
            // Empty list on error
        }

        return Inertia::render('Files/Index', [
            'websites' => $websites,
            'selectedWebsite' => $selectedWebsite,
            'currentPath' => $relativePath,
            'basePath' => $basePath,
            'files' => $files,
        ]);
    }

    public function read(Request $request): JsonResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
        ]);

        try {
            $content = $this->agentClient->readFile($validated['base_path'], $validated['relative_path']);
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

    public function write(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
            'content' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        try {
            $this->agentClient->writeFile($validated['base_path'], $validated['relative_path'], $validated['content']);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'file.write',
                'resource_type' => 'file',
                'resource_id' => $validated['relative_path'],
                'status' => 'success',
            ]);

            $url = route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ]);
            return redirect($url)->with('success', "File saved successfully.");
        } catch (Exception $e) {
            return back()->with('error', "Failed to write file: " . $e->getMessage());
        }
    }

    public function delete(Request $request): RedirectResponse
    {
        $validated = $request->validate([
            'base_path' => ['required', 'string'],
            'relative_path' => ['required', 'string'],
            'website_id' => ['nullable'],
            'current_path' => ['nullable', 'string'],
        ]);

        try {
            $this->agentClient->deleteFile($validated['base_path'], $validated['relative_path']);

            ActivityLog::create([
                'user_id' => $request->user()->id,
                'user_email' => $request->user()->email,
                'ip_address' => $request->ip() ?: '127.0.0.1',
                'user_agent' => $request->userAgent(),
                'action' => 'file.delete',
                'resource_type' => 'file',
                'resource_id' => $validated['relative_path'],
                'status' => 'success',
            ]);

            $url = route('files.index', [
                'website_id' => $validated['website_id'] ?? null,
                'path' => $validated['current_path'] ?? '',
            ]);
            return redirect($url)->with('success', "Item deleted successfully.");
        } catch (Exception $e) {
            return back()->with('error', "Failed to delete item: " . $e->getMessage());
        }
    }
}
