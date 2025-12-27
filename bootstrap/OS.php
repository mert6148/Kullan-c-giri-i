<?php

namespace Bootstrap;

use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Redis;
use Illuminate\Support\Facades\Mail;
use Illuminate\Support\Facades\Queue;
use Illuminate\Support\Facades\Session;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\Facades\Hash;

class OS {
    public function __construct() {
        $this->os = $this->getOS();
    }

    public function getOS() {
        return php_uname('s');
    }

    public function getCPU() {
        return php_uname('p');
    }

    public function getRAM() {
        return php_uname('m');
    }

    public function getDisk() {
        return php_uname('d');
    }

    public function getNetwork() {
        return php_uname('n');
    }

    public function getIP() {
        return php_uname('i');
    }

    public function getPython() {
        return php_uname('v');
    }

    public function getPHP() {
        return php_uname('r');
    }
    
    public function getNode() {
        return php_uname('N');
    }

    public function getJava() {
        return php_uname('J');
    }

    public function getAllInfo() {
        return [
            'os' => $this->getOS(),
            'cpu' => $this->getCPU(),
            'ram' => $this->getRAM(),
            'disk' => $this->getDisk(),
            'network' => $this->getNetwork(),
            'ip' => $this->getIP(),
            'python' => $this->getPython(),
            'php' => $this->getPHP(),
            'node' => $this->getNode(),
            'java' => $this->getJava()
        ];
    }
}

$os = new OS();
echo $os->getCPU();
echo $os->getRAM();
echo $os->getDisk();
echo $os->getOS();
echo $os->getNetwork();
echo $os->getIP();
echo $os->getPython();
echo $os->getPHP();
echo $os->getNode();
echo $os->getJava();
echo "\n";
echo json_encode($os->getAllInfo());