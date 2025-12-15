'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { register } from '@/lib/api';

export default function RegisterForm() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    name: '',
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  // 处理表单输入变化
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
    setError(''); // 清除错误信息
  };

  // 表单验证
  const validateForm = () => {
    if (!formData.email || !formData.password) {
      setError('邮箱和密码不能为空');
      return false;
    }

    // 邮箱格式验证
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(formData.email)) {
      setError('邮箱格式不正确');
      return false;
    }

    // 密码长度验证
    if (formData.password.length < 6) {
      setError('密码必须至少 6 个字符');
      return false;
    }

    // 确认密码验证
    if (formData.password !== formData.confirmPassword) {
      setError('两次输入的密码不一致');
      return false;
    }

    return true;
  };

  // 处理表单提交
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await register(formData.email, formData.password, formData.name);
      
      if (response.success) {
        // 注册成功，跳转到登录页
        router.push('/login?registered=true');
      } else {
        setError(response.error || '注册失败，请稍后重试');
      }
    } catch (err: any) {
      console.error('注册错误:', err);
      setError(err.response?.data?.error || err.response?.data?.message || '注册失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* 错误提示 */}
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
          {error}
        </div>
      )}

      {/* 用户名输入框 */}
      <div>
        <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-2">
          用户名 (可选)
        </label>
        <input
          type="text"
          id="name"
          name="name"
          value={formData.name}
          onChange={handleChange}
          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
          placeholder="输入您的用户名"
        />
      </div>

      {/* 邮箱输入框 */}
      <div>
        <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
          邮箱 <span className="text-red-500">*</span>
        </label>
        <input
          type="email"
          id="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          required
          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
          placeholder="your@email.com"
        />
      </div>

      {/* 密码输入框 */}
      <div>
        <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
          密码 <span className="text-red-500">*</span>
        </label>
        <input
          type="password"
          id="password"
          name="password"
          value={formData.password}
          onChange={handleChange}
          required
          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
          placeholder="至少 6 个字符"
        />
      </div>

      {/* 确认密码输入框 */}
      <div>
        <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-2">
          确认密码 <span className="text-red-500">*</span>
        </label>
        <input
          type="password"
          id="confirmPassword"
          name="confirmPassword"
          value={formData.confirmPassword}
          onChange={handleChange}
          required
          className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-all"
          placeholder="再次输入密码"
        />
      </div>

      {/* 提交按钮 */}
      <button
        type="submit"
        disabled={loading}
        className="w-full bg-indigo-600 text-white py-3 px-4 rounded-lg font-medium hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? '注册中...' : '注册'}
      </button>
    </form>
  );
}

