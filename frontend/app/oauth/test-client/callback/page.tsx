'use client';

import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

export default function CallbackPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    // 获取授权码
    const code = searchParams.get('code');
    const error = searchParams.get('error');
    const state = searchParams.get('state');

    if (error) {
      // 授权失败，跳转回测试客户端并显示错误
      router.push(`/oauth/test-client?error=${error}`);
      return;
    }

    if (code) {
      // 授权成功，跳转回测试客户端并传递授权码
      router.push(`/oauth/test-client?code=${code}&state=${state}`);
    }
  }, [router, searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600 mx-auto mb-4"></div>
        <p className="text-gray-600">处理授权回调中...</p>
      </div>
    </div>
  );
}

