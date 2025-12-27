<?php

namespace App\Cache;

use Symfony\Component\Cache\Adapter\FilesystemAdapter;
use Symfony\Component\Cache\Traits\FilesystemProxy;
use Symfony\Component\Cache\Traits\RedisProxy;
use Symfony\Component\Cache\Traits\MemcacheIProxy;

class Login
{
    public function __construct()
    {
        $this->cache = new Cache();
    }
}

?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login System</title>
    <style type="text/css" inline>
        .login-container {
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            padding: 20px;
            background: linear-gradient(135deg, #0a2a09 0%, #0f3d0b 50%, #1a5a15 100%);
        }

        .login-box {
            background-color: var(--bg-secondary);
            border: 2px solid var(--border-color);
            border-radius: 15px;
            padding: 40px;
            width: 100%;
            max-width: 450px;
            box-shadow: var(--shadow), 0 0 20px rgba(135, 194, 18, 0.2);
            animation: fadeIn 0.5s ease-in;
        }

        .login-header {
            text-align: center;
            margin-bottom: 30px;
        }
    </style>

    <script type="text/javascript" inline>
        function login() {
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const loginForm = document.getElementById('login-form');
            loginForm.addEventListener('submit', (e) => {
                e.preventDefault();
                const username = document.getElementById('username').value;
                const password = document.getElementById('password').value;
                console.log(username, password);
            });
        }

        const loginForm = document.getElementById('login-form');
        loginForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            console.log(username, password);
        });

        function logout() {
            const loginForm = document.getElementById('login-form');
            loginForm.addEventListener('submit', (e) => {
                e.preventDefault();
                const username = document.getElementById('username').value;
                const password = document.getElementById('password').value;
                console.log(username, password);
            });
        }

        function register() {
            const loginForm = document.getElementById('login-form');
            loginForm.addEventListener('submit', (e) => {
                e.preventDefault();
                const username = document.getElementById('username').value;
                const password = document.getElementById('password').value;
                console.log(username, password);
            });
        }
    </script>
</head>
<body>
    <div class="login-container">
        <div class="login-box">
            <div class="login-header">
                <h1>Login System</h1>
            </div>
            <div class="login-form">
                <input type="text" name="username" placeholder="Username">
                <input type="password" name="password" placeholder="Password">
                <button type="submit">Login</button>
            </div>
        </div>

        <div class="login-footer">
            <p>Copyright 2025 Login System</p>
        </div>

        <header>
            <h1>Login System</h1>
            <button onclick="login()">Login</button>
            <button onclick="logout()">Logout</button>
            <button onclick="register()">Register</button>
            <button onclick="forgotPassword()">Forgot Password</button>
        </header>

        <footer>
            <p>Copyright 2025 Login System</p>
            <button onclick="login()">Login</button>
            <button onclick="logout()">Logout</button>
            <button onclick="register()">Register</button>
            <button onclick="forgotPassword()">Forgot Password</button>
        </footer>

        <aside>
            <p>Copyright 2025 Login System</p>
            <button onclick="login()">Login</button>
            <button onclick="logout()">Logout</button>
            <button onclick="register()">Register</button>
            <button onclick="forgotPassword()">Forgot Password</button>
        </aside>
    </div>

    <div class="copyright">
        <div idate="container">

        </div>
    </div>
</body>
</html>