<?php

namespace Config\Routes;

use Illuminate\Support\Facades\Route;
use App\Controllers\HomeController;
use App\Controllers\AboutController;
use App\Controllers\ContactController;
use App\Controllers\LoginController;
use App\Controllers\RegisterController;
use App\Controllers\DashboardController;
use App\Controllers\ProfileController;
use App\Controllers\SettingsController;
use App\Controllers\LogoutController;

Route::get('/', 'HomeController@index')->name('home');
Route::get('/about', 'AboutController@index')->name('about');
Route::get('/contact', 'ContactController@index')->name('contact');
Route::get('/login', 'LoginController@index')->name('login');
Route::get('/register', 'RegisterController@index')->name('register');
Route::get('/dashboard', 'DashboardController@index')->name('dashboard');
Route::get('/profile', 'ProfileController@index')->name('profile');
Route::get('/settings', 'SettingsController@index')->name('settings');
Route::get('/logout', 'LogoutController@index')->name('logout');
Route::get('/admin', 'AdminController@index')->name('admin');
Route::get('/admin/login', 'AdminLoginController@index')->name('admin.login');
Route::get('/admin/logout', 'AdminLogoutController@index')->name('admin.logout');
Route::get('/admin/dashboard', 'AdminDashboardController@index')->name('admin.dashboard');
Route::get('/admin/profile', 'AdminProfileController@index')->name('admin.profile');
Route::get('/admin/settings', 'AdminSettingsController@index')->name('admin.settings');
Route::get('/admin/logout', 'AdminLogoutController@index')->name('admin.logout');

if (env('APP_ENV') === 'local') {
    Route::get('/test', function () {
        return view('test');
    })->name('test');
    Route::get('/test/login', function () {
        return view('test.login');
    })->name('test.login');
    Route::get('/test/register', function () {
        return view('test.register');
    })->name('test.register');
    Route::get('/test/dashboard', function () {
        return view('test.dashboard');
    })->name('test.dashboard');
    Route::get('/test/profile', function () {
        return view('test.profile');
    })->name('test.profile');
    Route::get('/test/settings', function () {
        return view('test.settings');
    })->name('test.settings');
    Route::get('/test/logout', function () {
        return view('test.logout');
    })->name('test.logout');
    Route::get('/test/admin', function () {
        return view('test.admin');
    })->name('test.admin');
    Route::get('/test/admin/login', function () {
        return view('test.admin.login');
    })->name('test.admin.login');
    Route::get('/test/admin/logout', function () {
        return view('test.admin.logout');
    })->name('test.admin.logout');
}

class HomeController extends Controller {
    public function index() {
        return view('layouts.app', [
            'title' => 'Home',
            'description' => 'Home page',
            'keywords' => 'home, page',
            'author' => 'Home',
            'copyright' => 'Home',
            'contact' => 'Home',
            'contact_email' => 'Home',
            'contact_phone' => 'Home',
            'contact_address' => 'Home',
            'contact_city' => 'Home',
            'contact_state' => 'Home',
            'contact_zip' => 'Home',
            'contact_country' => 'Home',
            'contact_website' => 'Home',
            'contact_facebook' => 'Home',
            'contact_twitter' => 'Home',
            'contact_instagram' => 'Home',
            'contact_linkedin' => 'Home',
            'contact_youtube' => 'Home',
            'contact_pinterest' => 'Home',
            'contact_tiktok' => 'Home',
            'contact_snapchat' => 'Home',
            'contact_reddit' => 'Home',
            'contact_github' => 'Home',
            'contact_gitlab' => 'Home',
            'contact_bitbucket' => 'Home'
        ]);
    }

    public function about() {
        return view('about');
    }

    public function contact() {
        return view('contact');
    }

    public function login() {
        return view('login');
    }

    public function register() {
        return view('register');
    }

    public function dashboard() {
        return view('dashboard');
    }
    
    public function profile() {
        return view('profile');
    }

    public function settings() {
        return view('settings');
    }

    public function logout() {
        return view('logout');
    }
}

class AboutController extends Controller {
    public function about() {
        return view('about');
    }

    public function contact() {
        return view('contact');
    }

    public function login() {
        return view('login');
    }

    public function register() {
        return view('register');
    }

    public function dashboard() {
        return view('dashboard');
    }
    
    public function profile() {
        return view('profile');
    }

    public function settings() {
        return view('settings');
    }

    public function logout() {
        return view('logout');
    }
}

class ContactController extends Controller {
    public function about() {
        return view('about');
    }

    public function contact() {
        return view('contact');
    }

    public function login() {
        return view('login');
    }

    public function register() {
        return view('register');
    }

    public function dashboard() {
        return view('dashboard');
    }
    
    public function profile() {
        return view('profile');
    }

    public function settings() {
        return view('settings');
    }

    public function logout() {
        return view('logout');
    }
}