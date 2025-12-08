#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Go Job 博客系统 API 测试脚本
测试所有接口的功能和边界情况
"""

import requests
import json
import time
from datetime import datetime

class APITester:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.token = None
        self.user_id = None
        self.test_results = []

    def log_test(self, test_name, method, url, payload=None, expected_status=200, actual_status=None, response=None, success=False):
        """记录测试结果"""
        result = {
            "test_name": test_name,
            "method": method,
            "url": url,
            "payload": payload,
            "expected_status": expected_status,
            "actual_status": actual_status,
            "response": response,
            "success": success,
            "timestamp": datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        }
        self.test_results.append(result)

        status_icon = "[PASS]" if success else "[FAIL]"
        print(f"{status_icon} {test_name}")
        if not success:
            print(f"   期望状态码: {expected_status}, 实际状态码: {actual_status}")
            print(f"   响应: {response}")
        print()

    def make_request(self, method, endpoint, payload=None, headers=None, expected_status=200):
        """发送HTTP请求"""
        url = f"{self.base_url}{endpoint}"

        if headers is None:
            headers = {"Content-Type": "application/json"}

        if self.token and "Authorization" not in headers:
            headers["Authorization"] = f"Bearer {self.token}"

        try:
            if method.upper() == "GET":
                response = requests.get(url, headers=headers, params=payload)
            elif method.upper() == "POST":
                response = requests.post(url, json=payload, headers=headers)
            elif method.upper() == "PUT":
                response = requests.put(url, json=payload, headers=headers)
            elif method.upper() == "DELETE":
                response = requests.delete(url, headers=headers)

            return response
        except requests.exceptions.RequestException as e:
            print(f"请求失败: {e}")
            return None

    def test_user_registration(self):
        """测试用户注册"""
        print("测试用户认证相关接口")
        print("=" * 50)

        # 测试1: 正常注册
        payload = {
            "username": "testuser",
            "email": "test@example.com",
            "password": "123456"
        }

        response = self.make_request("POST", "/auth/register", payload, expected_status=201)
        if response:
            success = response.status_code == 201
            response_data = response.json() if response.content else {}

            if success and "token" in response_data:
                self.token = response_data["token"]
                if "user" in response_data and "id" in response_data["user"]:
                    self.user_id = response_data["user"]["id"]

            self.log_test(
                "用户注册 - 正常情况",
                "POST", "/auth/register", payload, 201,
                response.status_code, response_data, success
            )

        # 测试2: 重复注册（应该失败）
        response = self.make_request("POST", "/auth/register", payload, expected_status=400)
        if response:
            success = response.status_code == 400
            response_data = response.json() if response.content else {}

            self.log_test(
                "用户注册 - 重复邮箱",
                "POST", "/auth/register", payload, 400,
                response.status_code, response_data, success
            )

        # 测试3: 无效邮箱格式
        invalid_payload = {
            "username": "testuser2",
            "email": "invalid-email",
            "password": "123456"
        }

        response = self.make_request("POST", "/auth/register", invalid_payload, expected_status=400)
        if response:
            success = response.status_code == 400
            response_data = response.json() if response.content else {}

            self.log_test(
                "用户注册 - 无效邮箱格式",
                "POST", "/auth/register", invalid_payload, 400,
                response.status_code, response_data, success
            )

        # 测试4: 密码太短
        short_password_payload = {
            "username": "testuser3",
            "email": "test3@example.com",
            "password": "123"
        }

        response = self.make_request("POST", "/auth/register", short_password_payload, expected_status=400)
        if response:
            success = response.status_code == 400
            response_data = response.json() if response.content else {}

            self.log_test(
                "用户注册 - 密码太短",
                "POST", "/auth/register", short_password_payload, 400,
                response.status_code, response_data, success
            )

    def test_user_login(self):
        """测试用户登录"""
        # 测试1: 正常登录
        payload = {
            "email": "test@example.com",
            "password": "123456"
        }

        response = self.make_request("POST", "/auth/login", payload, expected_status=200)
        if response:
            success = response.status_code == 200
            response_data = response.json() if response.content else {}

            if success and "token" in response_data:
                self.token = response_data["token"]

            self.log_test(
                "用户登录 - 正常情况",
                "POST", "/auth/login", payload, 200,
                response.status_code, response_data, success
            )

        # 测试2: 错误密码
        wrong_password_payload = {
            "email": "test@example.com",
            "password": "wrongpassword"
        }

        response = self.make_request("POST", "/auth/login", wrong_password_payload, expected_status=401)
        if response:
            success = response.status_code == 401
            response_data = response.json() if response.content else {}

            self.log_test(
                "用户登录 - 错误密码",
                "POST", "/auth/login", wrong_password_payload, 401,
                response.status_code, response_data, success
            )

        # 测试3: 不存在的用户
        nonexistent_user_payload = {
            "email": "nonexistent@example.com",
            "password": "123456"
        }

        response = self.make_request("POST", "/auth/login", nonexistent_user_payload, expected_status=401)
        if response:
            success = response.status_code == 401
            response_data = response.json() if response.content else {}

            self.log_test(
                "用户登录 - 不存在的用户",
                "POST", "/auth/login", nonexistent_user_payload, 401,
                response.status_code, response_data, success
            )

    def test_token_refresh(self):
        """测试Token刷新"""
        if not self.token:
            print("[WARNING] 跳过Token刷新测试 - 没有有效的Token")
            return

        headers = {"Authorization": f"Bearer {self.token}"}
        response = self.make_request("POST", "/auth/refresh", headers=headers, expected_status=200)

        if response:
            success = response.status_code == 200
            response_data = response.json() if response.content else {}

            if success and "token" in response_data:
                self.token = response_data["token"]

            self.log_test(
                "Token刷新 - 正常情况",
                "POST", "/auth/refresh", None, 200,
                response.status_code, response_data, success
            )

    def test_user_profile(self):
        """测试获取用户信息"""
        print("测试用户信息相关接口")
        print("=" * 50)

        if not self.token:
            print("[WARNING] 跳过用户信息测试 - 没有有效的Token")
            return

        # 测试1: 正常获取用户信息
        response = self.make_request("GET", "/api/profile", expected_status=200)
        if response:
            success = response.status_code == 200
            response_data = response.json() if response.content else {}

            self.log_test(
                "获取用户信息 - 正常情况",
                "GET", "/api/profile", None, 200,
                response.status_code, response_data, success
            )

        # 测试2: 无Token访问
        # 临时保存token并清空，确保不会自动添加Authorization header
        temp_token = self.token
        self.token = None
        headers = {"Content-Type": "application/json"}
        response = self.make_request("GET", "/api/profile", headers=headers, expected_status=401)
        # 恢复token
        self.token = temp_token
        if response:
            success = response.status_code == 401
            response_data = response.json() if response.content else {}

            self.log_test(
                "获取用户信息 - 无Token",
                "GET", "/api/profile", None, 401,
                response.status_code, response_data, success
            )

    def test_posts(self):
        """测试文章相关接口"""
        print("测试文章相关接口")
        print("=" * 50)

        if not self.token:
            print("[WARNING] 跳过文章测试 - 没有有效的Token")
            return

        # 测试1: 创建文章
        post_payload = {
            "title": "测试文章标题",
            "content": "这是一篇测试文章的内容"
        }

        response = self.make_request("POST", "/api/create_post", post_payload, expected_status=201)
        if response:
            success = response.status_code == 201
            response_data = response.json() if response.content else {}

            self.log_test(
                "创建文章 - 正常情况",
                "POST", "/api/create_post", post_payload, 201,
                response.status_code, response_data, success
            )

        # 测试2: 获取文章列表
        response = self.make_request("GET", "/api/get_post", expected_status=200)
        if response:
            success = response.status_code == 200
            response_data = response.json() if response.content else {}

            self.log_test(
                "获取文章列表 - 正常情况",
                "GET", "/api/get_post", None, 200,
                response.status_code, response_data, success
            )

        # 测试3: 获取特定文章
        response = self.make_request("GET", "/api/get_post?id=1", expected_status=200)
        if response:
            success = response.status_code == 200
            response_data = response.json() if response.content else {}

            self.log_test(
                "获取特定文章 - 正常情况",
                "GET", "/api/get_post?id=1", None, 200,
                response.status_code, response_data, success
            )

        # 测试4: 更新文章
        update_payload = {
            "title": "更新后的文章标题",
            "content": "更新后的文章内容"
        }

        response = self.make_request("POST", "/api/update_post?id=1&title=更新后的文章标题&content=更新后的文章内容", expected_status=200)
        if response:
            success = response.status_code == 200
            response_data = response.json() if response.content else {}

            self.log_test(
                "更新文章 - 正常情况",
                "POST", "/api/update_post?id=1", update_payload, 200,
                response.status_code, response_data, success
            )

    def test_comments(self):
        """测试评论相关接口"""
        print("测试评论相关接口")
        print("=" * 50)

        if not self.token:
            print("[WARNING] 跳过评论测试 - 没有有效的Token")
            return

        # 测试1: 创建评论
        comment_payload = {
            "content": "这是一条测试评论",
            "post_id": 1
        }

        response = self.make_request("POST", "/api/create_comment", comment_payload, expected_status=201)
        if response:
            success = response.status_code == 201
            response_data = response.json() if response.content else {}

            self.log_test(
                "创建评论 - 正常情况",
                "POST", "/api/create_comment", comment_payload, 201,
                response.status_code, response_data, success
            )

        # 测试2: 获取评论
        response = self.make_request("GET", "/api/get_comment?id=1", expected_status=200)
        if response:
            success = response.status_code == 200
            response_data = response.json() if response.content else {}

            self.log_test(
                "获取评论 - 正常情况",
                "GET", "/api/get_comment?id=1", None, 200,
                response.status_code, response_data, success
            )

    def test_public_endpoints(self):
        """测试公开接口"""
        print("测试公开接口")
        print("=" * 50)

        # 测试1: 公开测试接口（无Token）
        headers = {"Content-Type": "application/json"}
        response = self.make_request("GET", "/public/test", headers=headers, expected_status=200)
        if response:
            success = response.status_code == 200
            response_data = response.json() if response.content else {}

            self.log_test(
                "公开测试接口 - 无Token",
                "GET", "/public/test", None, 200,
                response.status_code, response_data, success
            )

        # 测试2: 公开测试接口（有Token）
        if self.token:
            response = self.make_request("GET", "/public/test", expected_status=200)
            if response:
                success = response.status_code == 200
                response_data = response.json() if response.content else {}

                self.log_test(
                    "公开测试接口 - 有Token",
                    "GET", "/public/test", None, 200,
                    response.status_code, response_data, success
                )

    def generate_report(self):
        """生成测试报告"""
        print("\n" + "=" * 60)
        print("测试报告")
        print("=" * 60)

        total_tests = len(self.test_results)
        passed_tests = sum(1 for result in self.test_results if result["success"])
        failed_tests = total_tests - passed_tests

        print(f"总测试数: {total_tests}")
        print(f"通过: {passed_tests} [PASS]")
        print(f"失败: {failed_tests} [FAIL]")
        print(f"成功率: {(passed_tests/total_tests*100):.1f}%")

        if failed_tests > 0:
            print("\n[FAIL] 失败的测试:")
            for result in self.test_results:
                if not result["success"]:
                    print(f"  - {result['test_name']}")
                    print(f"    期望: {result['expected_status']}, 实际: {result['actual_status']}")
                    print(f"    响应: {result['response']}")

        # 保存详细报告到文件
        with open("api_test_report.json", "w", encoding="utf-8") as f:
            json.dump(self.test_results, f, ensure_ascii=False, indent=2)

        print(f"\n详细测试报告已保存到: api_test_report.json")

    def run_all_tests(self):
        """运行所有测试"""
        print("开始API接口测试")
        print("=" * 60)
        print(f"测试目标: {self.base_url}")
        print(f"开始时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print()

        # 检查服务器是否可用
        try:
            response = requests.get(f"{self.base_url}/public/test", timeout=5)
            if response.status_code != 200:
                print("[ERROR] 服务器不可用，请确保服务器正在运行")
                return
        except requests.exceptions.RequestException:
            print("[ERROR] 无法连接到服务器，请确保服务器正在运行")
            return

        print("[OK] 服务器连接正常\n")

        # 运行测试
        self.test_user_registration()
        self.test_user_login()
        self.test_token_refresh()
        self.test_user_profile()
        self.test_posts()
        self.test_comments()
        self.test_public_endpoints()

        # 生成报告
        self.generate_report()

if __name__ == "__main__":
    # 创建测试实例并运行测试
    tester = APITester()
    tester.run_all_tests()
