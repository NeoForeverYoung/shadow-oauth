import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="max-w-md w-full space-y-8 p-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            🔐 Shadow IAM
          </h1>
          <p className="text-lg text-gray-600 mb-8">
            现代化身份认证与访问管理系统
          </p>
        </div>

        <div className="bg-white rounded-lg shadow-xl p-8 space-y-6">
          <div className="space-y-4">
            <Link
              href="/login"
              className="block w-full bg-indigo-600 text-white text-center py-3 px-4 rounded-lg font-medium hover:bg-indigo-700 transition-colors"
            >
              登录
            </Link>
            <Link
              href="/register"
              className="block w-full bg-white text-indigo-600 text-center py-3 px-4 rounded-lg font-medium border-2 border-indigo-600 hover:bg-indigo-50 transition-colors"
            >
              注册新账户
            </Link>
          </div>

          <div className="pt-4 border-t border-gray-200">
            <h3 className="text-sm font-semibold text-gray-700 mb-2">✨ 核心功能</h3>
            <ul className="text-sm text-gray-600 space-y-1">
              <li>• 安全的用户注册和登录</li>
              <li>• JWT Token 认证</li>
              <li>• 用户信息管理</li>
              <li>• 密码加密保护</li>
            </ul>
          </div>
        </div>

        <p className="text-center text-sm text-gray-500">
          参考 Casdoor 设计 | Golang + Gin + Next.js
        </p>
      </div>
    </div>
  );
}

