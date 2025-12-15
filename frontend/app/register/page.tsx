import Link from 'next/link';
import RegisterForm from '@/components/RegisterForm';

export default function RegisterPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {/* 顶部标题 */}
        <div className="text-center">
          <h2 className="text-3xl font-bold text-gray-900">
            创建新账户
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            已有账户？{' '}
            <Link href="/login" className="font-medium text-indigo-600 hover:text-indigo-500">
              立即登录
            </Link>
          </p>
        </div>

        {/* 注册表单卡片 */}
        <div className="bg-white rounded-lg shadow-xl p-8">
          <RegisterForm />
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

