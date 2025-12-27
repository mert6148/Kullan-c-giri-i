<?php

use Illuminate\Support\Facades\Route;
use App\Http\Controllers\UserManageController;

Route::prefix('user-manage')->name('user-manage.')->group(function () {
    Route::get('/', [UserManageController::class, 'index'])->name('index');
    Route::get('/create', [UserManageController::class, 'create'])->name('create');
    Route::get('/edit/{id}', [UserManageController::class, 'edit'])->name('edit');
    Route::delete('/delete/{id}', [UserManageController::class, 'destroy'])->name('delete');
    Route::get('/login', [UserManageController::class, 'login'])->name('login');
    Route::get('/register', [UserManageController::class, 'register'])->name('register');
    Route::get('/forgot-password', [UserManageController::class, 'forgotPassword'])->name('forgot-password');
    Route::get('/reset-password/{token}', [UserManageController::class, 'resetPassword'])->name('reset-password');
});

Route::post('/user-manage/store', [UserManageController::class, 'store'])->name('user-manage.store');
Route::put('/user-manage/update/{id}', [UserManageController::class, 'update'])->name('user-manage.update');
Route::post('/user-manage/authenticate', [UserManageController::class, 'authenticate'])->name('user-manage.authenticate');
Route::post('/user-manage/logout', [UserManageController::class, 'logout'])->name('user-manage.logout');

?>