<?php

namespace Config\Routes;

use Illuminate\Support\Facades\Route;

class Services extends Controller {
    public function index() {
        return view('services');
    }
    public function create() {
        return view('services.create');
    }
    public function edit($id) {
        return view('services.edit', ['id' => $id]);
    }
    public function delete($id) {
        return view('services.delete', ['id' => $id]);
    }

    public function show($id) {
        return view('services.show', ['id' => $id]);
    }

    public function store($id) {
        return view('services.store');
    }

    public function update($id) {
        return view('services.update', ['id' => $id]);
    }
}

class Config extends Routes {
    public function update($id) {
        return view('services.update', ['id' => $id]);
    }

    public function strtolower() {
        $retVal = (condition) ? a : b ;
    }
   
    public function strtoupper() {
        $retVal = (condition) ? a : b ;
    }

}