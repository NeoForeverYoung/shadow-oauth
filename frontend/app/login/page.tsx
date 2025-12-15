'use client';

import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import LoginForm from '@/components/LoginForm';
import { Suspense } from 'react';

function LoginPageContent() {
  const searchParams = useSearchParams();
  const registered = searchParams.get('registered');

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {/* 顶部标题 */}
        <div className="text-center">
          <h2 className="text-3xl font-bold text-gray-900">
            登录您的账户
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            还没有账户？{' '}
            <Link href="/register" className="font-medium text-indigo-600 hover:text-indigo-500">
              立即注册
            </Link>
          </p>
        </div>

        {/* 注册成功提示 */}
        {registered === 'true' && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg text-sm text-center">
            ✅ 注册成功！请使用您的邮箱和密码登录
          </div>
        )}

        {/* 登录表单卡片 */}
        <div className="bg-white rounded-lg shadow-xl p-8">
          <LoginForm />
        </div>

        {/* 返回首页链接 */}
        <div className="text-center">
          <Link href="/" className="text-sm text-gray-600 hover:text-gray-900">
            ← 返回首页
          </Link>
        </div>
      </div>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="text-gray-600">加载中...</div>
      </div>
    }>
      <LoginPageContent />
    </Suspense>
  );
}

