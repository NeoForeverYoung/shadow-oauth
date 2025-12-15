'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getCurrentUser, User } from '@/lib/api';
import { isAuthenticated, logout } from '@/lib/auth';

export default function DashboardPage() {
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // 检查认证状态并获取用户信息
  useEffect(() => {
    const fetchUserData = async () => {
      // 检查是否已登录
      if (!isAuthenticated()) {
        router.push('/login');
        return;
      }

      try {
        // 从后端获取当前用户信息
        const response = await getCurrentUser();
        
        if (response.success && response.data) {
          setUser(response.data);
        } else {
          setError('获取用户信息失败');
          setTimeout(() => {
            logout();
          }, 2000);
        }
      } catch (err: any) {
        console.error('获取用户信息错误:', err);
        setError('获取用户信息失败，请重新登录');
        setTimeout(() => {
          logout();
        }, 2000);
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, [router]);

  // 处理登出
  const handleLogout = () => {
    logout();
  };

  // 加载中状态
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

  // 错误状态
  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="bg-white rounded-lg shadow-xl p-8 max-w-md">
          <div className="text-center">
            <div className="text-red-500 text-5xl mb-4">⚠️</div>
            <h2 className="text-xl font-bold text-gray-900 mb-2">出错了</h2>
            <p className="text-gray-600">{error}</p>
          </div>
        </div>
      </div>
    );
  }

  // 正常显示仪表盘
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* 顶部导航栏 */}
      <nav className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900">🔐 Shadow IAM</h1>
            </div>
            <button
              onClick={handleLogout}
              className="bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600 transition-colors text-sm font-medium"
            >
              退出登录
            </button>
          </div>
        </div>
      </nav>

      {/* 主内容区域 */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* 欢迎卡片 */}
        <div className="bg-white rounded-lg shadow-xl p-8 mb-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">
            欢迎回来，{user?.name || '用户'}！👋
          </h2>
          <p className="text-gray-600">
            您已成功登录 Shadow IAM 系统
          </p>
        </div>

        {/* 用户信息卡片 */}
        <div className="bg-white rounded-lg shadow-xl p-8 mb-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">📋 个人信息</h3>
          <div className="space-y-3">
            <div className="flex items-center">
              <span className="text-gray-600 w-32">用户 ID:</span>
              <span className="font-medium text-gray-900">{user?.id}</span>
            </div>
            <div className="flex items-center">
              <span className="text-gray-600 w-32">用户名:</span>
              <span className="font-medium text-gray-900">{user?.name || '未设置'}</span>
            </div>
            <div className="flex items-center">
              <span className="text-gray-600 w-32">邮箱:</span>
              <span className="font-medium text-gray-900">{user?.email}</span>
            </div>
            <div className="flex items-center">
              <span className="text-gray-600 w-32">注册时间:</span>
              <span className="font-medium text-gray-900">
                {user?.created_at ? new Date(user.created_at).toLocaleString('zh-CN') : '-'}
              </span>
            </div>
            <div className="flex items-center">
              <span className="text-gray-600 w-32">最后更新:</span>
              <span className="font-medium text-gray-900">
                {user?.updated_at ? new Date(user.updated_at).toLocaleString('zh-CN') : '-'}
              </span>
            </div>
          </div>
        </div>

        {/* 功能状态卡片 */}
        <div className="bg-white rounded-lg shadow-xl p-8">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">✅ 已实现功能</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex items-start">
              <span className="text-green-500 mr-2">✓</span>
              <div>
                <p className="font-medium text-gray-900">用户注册</p>
                <p className="text-sm text-gray-600">邮箱密码注册，密码加密存储</p>
              </div>
            </div>
            <div className="flex items-start">
              <span className="text-green-500 mr-2">✓</span>
              <div>
                <p className="font-medium text-gray-900">用户登录</p>
                <p className="text-sm text-gray-600">JWT Token 认证机制</p>
              </div>
            </div>
            <div className="flex items-start">
              <span className="text-green-500 mr-2">✓</span>
              <div>
                <p className="font-medium text-gray-900">会话管理</p>
                <p className="text-sm text-gray-600">Token 自动续期和过期处理</p>
              </div>
            </div>
            <div className="flex items-start">
              <span className="text-green-500 mr-2">✓</span>
              <div>
                <p className="font-medium text-gray-900">路由保护</p>
                <p className="text-sm text-gray-600">未登录自动跳转登录页</p>
              </div>
            </div>
          </div>

          <div className="mt-6 pt-6 border-t border-gray-200">
            <h4 className="text-sm font-semibold text-gray-700 mb-2">🚀 后续扩展方向</h4>
            <ul className="text-sm text-gray-600 space-y-1">
              <li>• 用户角色和权限管理（RBAC）</li>
              <li>• OAuth 2.0 授权服务器</li>
              <li>• 组织/租户管理</li>
              <li>• 多因素认证（MFA）</li>
              <li>• 第三方登录（Google、GitHub等）</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}

