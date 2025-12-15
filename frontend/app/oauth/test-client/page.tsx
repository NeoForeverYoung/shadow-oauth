'use client';

import { useState, useEffect } from 'react';

// 测试客户端配置（简化版，实际应用中这些信息应该从服务器获取）
const TEST_CLIENT = {
  client_id: 'test_client_123',
  client_secret: 'test_secret_456',
  redirect_uri: 'http://localhost:3000/oauth/test-client/callback',
  name: 'OAuth 测试客户端',
};

export default function TestClientPage() {
  const [step, setStep] = useState<'start' | 'authorizing' | 'callback' | 'token' | 'userinfo'>('start');
  const [authCode, setAuthCode] = useState('');
  const [accessToken, setAccessToken] = useState('');
  const [userInfo, setUserInfo] = useState<any>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // 步骤 1: 开始 OAuth 流程
  const handleStartOAuth = () => {
    setError('');
    setStep('authorizing');

    // 构建授权 URL
    const authorizeUrl = new URL('/oauth/authorize', 'http://localhost:8080');
    authorizeUrl.searchParams.set('client_id', TEST_CLIENT.client_id);
    authorizeUrl.searchParams.set('redirect_uri', TEST_CLIENT.redirect_uri);
    authorizeUrl.searchParams.set('response_type', 'code');
    authorizeUrl.searchParams.set('state', 'test_state_123');

    // 跳转到授权页面（需要先登录）
    window.location.href = authorizeUrl.toString();
  };

  // 步骤 2: 处理回调（从 URL 获取授权码）
  useEffect(() => {
    if (typeof window !== 'undefined' && step === 'start') {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get('code');
      const error = urlParams.get('error');

      if (error) {
        setError(`授权失败: ${error}`);
        setStep('start');
        // 清除 URL 参数
        window.history.replaceState({}, '', '/oauth/test-client');
        return;
      }

      if (code) {
        setAuthCode(code);
        setStep('token');
        // 清除 URL 参数
        window.history.replaceState({}, '', '/oauth/test-client');
      }
    }
  }, [step]);

  // 步骤 3: 用授权码交换 Access Token
  const handleExchangeToken = async () => {
    if (!authCode) return;

    setLoading(true);
    setError('');

    try {
      // 构建 Token 请求（使用 form-urlencoded 格式）
      const formData = new URLSearchParams();
      formData.append('grant_type', 'authorization_code');
      formData.append('code', authCode);
      formData.append('redirect_uri', TEST_CLIENT.redirect_uri);
      formData.append('client_id', TEST_CLIENT.client_id);
      formData.append('client_secret', TEST_CLIENT.client_secret);

      const response = await fetch('http://localhost:8080/oauth/token', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData.toString(),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || '获取 Token 失败');
      }

      setAccessToken(data.access_token);
      setStep('userinfo');
    } catch (err: any) {
      setError(err.message || '交换 Token 失败');
    } finally {
      setLoading(false);
    }
  };

  // 步骤 4: 使用 Access Token 获取用户信息
  const handleGetUserInfo = async () => {
    if (!accessToken) return;

    setLoading(true);
    setError('');

    try {
      const response = await fetch(
        `http://localhost:8080/oauth/userinfo?access_token=${accessToken}`
      );

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || '获取用户信息失败');
      }

      setUserInfo(data.data);
    } catch (err: any) {
      setError(err.message || '获取用户信息失败');
    } finally {
      setLoading(false);
    }
  };

  // 检查是否是回调页面（从 URL 参数获取授权码）
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get('code');
      if (code && step === 'start') {
        setAuthCode(code);
        setStep('token');
        // 清除 URL 参数
        window.history.replaceState({}, '', '/oauth/test-client');
      }
    }
  }, [step]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-lg shadow-xl p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            🧪 OAuth 2.0 测试客户端
          </h1>
          <p className="text-gray-600 mb-8">
            这是一个演示 OAuth 2.0 授权码流程的测试客户端
          </p>

          {/* 客户端信息 */}
          <div className="bg-gray-50 rounded-lg p-4 mb-8">
            <h2 className="font-semibold text-gray-900 mb-2">客户端配置</h2>
            <div className="text-sm text-gray-600 space-y-1">
              <p><strong>Client ID:</strong> {TEST_CLIENT.client_id}</p>
              <p><strong>Redirect URI:</strong> {TEST_CLIENT.redirect_uri}</p>
            </div>
          </div>

          {/* 错误提示 */}
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
              {error}
            </div>
          )}

          {/* 步骤显示 */}
          <div className="space-y-6">
            {/* 步骤 1: 开始授权 */}
            {step === 'start' && (
              <div>
                <h3 className="text-xl font-semibold text-gray-900 mb-4">
                  步骤 1: 开始 OAuth 授权流程
                </h3>
                <p className="text-gray-600 mb-4">
                  点击下面的按钮，将跳转到授权服务器，请求用户授权。
                </p>
                <button
                  onClick={handleStartOAuth}
                  className="bg-indigo-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-indigo-700 transition-colors"
                >
                  开始授权
                </button>
              </div>
            )}

            {/* 步骤 2: 授权中 */}
            {step === 'authorizing' && (
              <div>
                <h3 className="text-xl font-semibold text-gray-900 mb-4">
                  步骤 2: 等待用户授权...
                </h3>
                <p className="text-gray-600">
                  请在新打开的页面中登录并授权。
                </p>
              </div>
            )}

            {/* 步骤 3: 交换 Token */}
            {step === 'token' && (
              <div>
                <h3 className="text-xl font-semibold text-gray-900 mb-4">
                  步骤 3: 用授权码交换 Access Token
                </h3>
                <div className="bg-gray-50 rounded-lg p-4 mb-4">
                  <p className="text-sm text-gray-600 mb-2">收到的授权码：</p>
                  <p className="font-mono text-sm text-gray-800 break-all">{authCode}</p>
                </div>
                <button
                  onClick={handleExchangeToken}
                  disabled={loading}
                  className="bg-indigo-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50"
                >
                  {loading ? '交换中...' : '交换 Access Token'}
                </button>
              </div>
            )}

            {/* 步骤 4: 获取用户信息 */}
            {step === 'userinfo' && (
              <div>
                <h3 className="text-xl font-semibold text-gray-900 mb-4">
                  步骤 4: 使用 Access Token 获取用户信息
                </h3>
                <div className="bg-gray-50 rounded-lg p-4 mb-4">
                  <p className="text-sm text-gray-600 mb-2">Access Token：</p>
                  <p className="font-mono text-xs text-gray-800 break-all">
                    {accessToken.substring(0, 50)}...
                  </p>
                </div>
                <button
                  onClick={handleGetUserInfo}
                  disabled={loading}
                  className="bg-indigo-600 text-white px-6 py-3 rounded-lg font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50 mb-4"
                >
                  {loading ? '获取中...' : '获取用户信息'}
                </button>

                {/* 显示用户信息 */}
                {userInfo && (
                  <div className="bg-green-50 border border-green-200 rounded-lg p-4 mt-4">
                    <h4 className="font-semibold text-gray-900 mb-2">✅ 成功获取用户信息：</h4>
                    <pre className="text-sm text-gray-700 overflow-auto">
                      {JSON.stringify(userInfo, null, 2)}
                    </pre>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* OAuth 流程说明 */}
          <div className="mt-8 pt-8 border-t border-gray-200">
            <h3 className="font-semibold text-gray-900 mb-4">OAuth 2.0 授权码流程说明</h3>
            <ol className="list-decimal list-inside space-y-2 text-sm text-gray-600">
              <li>客户端引导用户到授权服务器</li>
              <li>用户登录并授权</li>
              <li>授权服务器返回授权码（通过重定向）</li>
              <li>客户端用授权码交换 Access Token</li>
              <li>客户端使用 Access Token 访问受保护资源</li>
            </ol>
          </div>
        </div>
      </div>
    </div>
  );
}

