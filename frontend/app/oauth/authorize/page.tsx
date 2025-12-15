'use client';

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';

export default function AuthorizePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [clientInfo, setClientInfo] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // 从 URL 获取 OAuth 参数
  const clientId = searchParams.get('client_id');
  const redirectUri = searchParams.get('redirect_uri');
  const responseType = searchParams.get('response_type');
  const state = searchParams.get('state');

  useEffect(() => {
    // 验证参数
    if (!clientId || !redirectUri || responseType !== 'code') {
      setError('缺少必要的授权参数');
      setLoading(false);
      return;
    }

    // 这里可以调用 API 获取客户端信息（简化版，直接显示）
    setClientInfo({
      name: '测试应用',
      client_id: clientId,
    });
    setLoading(false);
  }, [clientId, redirectUri, responseType]);

  // 处理授权同意
  const handleApprove = () => {
    if (!clientId || !redirectUri) return;

    // 构建授权 URL（后端会自动处理授权码生成和重定向）
    const authorizeUrl = new URL('/oauth/authorize', 'http://localhost:8080');
    authorizeUrl.searchParams.set('client_id', clientId);
    authorizeUrl.searchParams.set('redirect_uri', redirectUri);
    authorizeUrl.searchParams.set('response_type', 'code');
    if (state) {
      authorizeUrl.searchParams.set('state', state);
    }

    // 跳转到后端授权端点（需要登录）
    window.location.href = authorizeUrl.toString();
  };

  // 处理拒绝授权
  const handleDeny = () => {
    if (!redirectUri) return;

    // 重定向回客户端，带上错误信息
    const denyUrl = new URL(redirectUri);
    denyUrl.searchParams.set('error', 'access_denied');
    if (state) {
      denyUrl.searchParams.set('state', state);
    }
    window.location.href = denyUrl.toString();
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
          <p className="text-gray-600">加载中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="bg-white rounded-lg shadow-xl p-8 max-w-md">
          <div className="text-center">
            <div className="text-red-500 text-5xl mb-4">⚠️</div>
            <h2 className="text-xl font-bold text-gray-900 mb-2">授权错误</h2>
            <p className="text-gray-600 mb-4">{error}</p>
            <Link href="/" className="text-indigo-600 hover:text-indigo-500">
              返回首页
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-lg shadow-xl p-8">
          <div className="text-center mb-6">
            <h2 className="text-2xl font-bold text-gray-900 mb-2">
              🔐 授权请求
            </h2>
            <p className="text-gray-600">
              一个应用想要访问您的账户
            </p>
          </div>

          {/* 客户端信息 */}
          <div className="bg-gray-50 rounded-lg p-4 mb-6">
            <div className="flex items-center mb-3">
              <div className="w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center mr-3">
                <span className="text-2xl">📱</span>
              </div>
              <div>
                <h3 className="font-semibold text-gray-900">
                  {clientInfo?.name || '未知应用'}
                </h3>
                <p className="text-sm text-gray-500">
                  Client ID: {clientId?.substring(0, 20)}...
                </p>
              </div>
            </div>
          </div>

          {/* 权限说明 */}
          <div className="mb-6">
            <h4 className="text-sm font-semibold text-gray-700 mb-2">
              此应用将能够：
            </h4>
            <ul className="text-sm text-gray-600 space-y-1">
              <li className="flex items-start">
                <span className="text-green-500 mr-2">✓</span>
                <span>查看您的基本信息（邮箱、用户名）</span>
              </li>
            </ul>
          </div>

          {/* 重定向地址 */}
          <div className="mb-6 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
            <p className="text-xs text-gray-600 mb-1">授权后将重定向到：</p>
            <p className="text-xs font-mono text-gray-800 break-all">
              {redirectUri}
            </p>
          </div>

          {/* 操作按钮 */}
          <div className="space-y-3">
            <button
              onClick={handleApprove}
              className="w-full bg-indigo-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-indigo-700 transition-colors"
            >
              同意授权
            </button>
            <button
              onClick={handleDeny}
              className="w-full bg-gray-200 text-gray-700 py-3 px-4 rounded-lg font-medium hover:bg-gray-300 transition-colors"
            >
              拒绝
            </button>
          </div>

          <div className="mt-6 pt-6 border-t border-gray-200 text-center">
            <p className="text-xs text-gray-500">
              您需要先登录才能授权
            </p>
            <Link href="/login" className="text-xs text-indigo-600 hover:text-indigo-500">
              前往登录
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}

