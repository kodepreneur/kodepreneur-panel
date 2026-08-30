<?php

namespace Tests\Feature;

use App\Services\Agent\AgentClientInterface;
use App\Services\Agent\HttpAgentClient;
use App\Services\Agent\MockAgentClient;
use Illuminate\Support\Facades\Http;
use Tests\TestCase;

class AgentClientTest extends TestCase
{
    public function test_mock_agent_returns_valid_metrics(): void
    {
        $mock = new MockAgentClient();
        $metrics = $mock->getSystemMetrics();

        $this->assertArrayHasKey('cpu', $metrics);
        $this->assertArrayHasKey('memory', $metrics);
        $this->assertArrayHasKey('disk', $metrics);
        $this->assertArrayHasKey('load_average', $metrics);
        $this->assertGreaterThan(0, $metrics['cpu']['cores']);
    }

    public function test_http_agent_client_signs_requests_with_hmac(): void
    {
        $secret = 'test-secret-key-123';
        $client = new HttpAgentClient('http://127.0.0.1:8765', $secret);

        Http::fake([
            'http://127.0.0.1:8765/api/v1/system/metrics' => Http::response([
                'success' => true,
                'data' => [
                    'cpu' => ['usage_percent' => 15.0, 'cores' => 4],
                    'memory' => ['total_bytes' => 2000, 'used_bytes' => 1000, 'free_bytes' => 1000, 'usage_percent' => 50.0],
                    'disk' => ['mount' => '/', 'total_bytes' => 5000, 'used_bytes' => 1000, 'free_bytes' => 4000, 'usage_percent' => 20.0],
                    'load_average' => ['load1' => 0.1, 'load5' => 0.2, 'load15' => 0.3],
                ],
            ], 200),
        ]);

        $metrics = $client->getSystemMetrics();

        $this->assertEquals(15.0, $metrics['cpu']['usage_percent']);

        Http::assertSent(function ($request) use ($secret) {
            $hasTimestamp = $request->hasHeader('X-Kodepreneur-Timestamp');
            $hasNonce = $request->hasHeader('X-Kodepreneur-Nonce');
            $hasSignature = $request->hasHeader('X-Kodepreneur-Signature');

            $timestamp = $request->header('X-Kodepreneur-Timestamp')[0];
            $nonce = $request->header('X-Kodepreneur-Nonce')[0];
            $sig = $request->header('X-Kodepreneur-Signature')[0];

            $bodyHash = hash('sha256', (string) $request->body());
            $expectedMsg = "GET\n/api/v1/system/metrics\n{$timestamp}\n{$nonce}\n{$bodyHash}";
            $expectedSig = hash_hmac('sha256', $expectedMsg, $secret);

            return $hasTimestamp && $hasNonce && $hasSignature && hash_equals($expectedSig, $sig);
        });
    }

    public function test_mock_agent_returns_valid_website_traffic(): void
    {
        $mock = new MockAgentClient();
        $traffic = $mock->getWebsiteTraffic('testsite.com', '24h');

        $this->assertEquals('testsite.com', $traffic['domain']);
        $this->assertEquals('24h', $traffic['period']);
        $this->assertGreaterThan(0, $traffic['total_requests']);
        $this->assertArrayHasKey('time_series', $traffic);
        $this->assertArrayHasKey('top_paths', $traffic);
        $this->assertArrayHasKey('recent_requests', $traffic);
    }

    public function test_http_agent_client_fetches_website_traffic(): void
    {
        $secret = 'test-secret-key-123';
        $client = new HttpAgentClient('http://127.0.0.1:8765', $secret);

        Http::fake([
            'http://127.0.0.1:8765/api/v1/websites/demo.com/traffic?period=7d' => Http::response([
                'success' => true,
                'data' => [
                    'domain' => 'demo.com',
                    'period' => '7d',
                    'total_requests' => 1500,
                ],
            ], 200),
        ]);

        $traffic = $client->getWebsiteTraffic('demo.com', '7d');

        $this->assertEquals('demo.com', $traffic['domain']);
        $this->assertEquals('7d', $traffic['period']);
        $this->assertEquals(1500, $traffic['total_requests']);
    }
}
