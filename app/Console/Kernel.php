<?php

namespace App;

use Kernel\Kernel;
use Kernel\Router;

if (session_status() == PHP_SESSION_NONE) {
   session_start();
   $kernel = new Kernel();
   $kernel->run();
}

public function getRouter(): Router
{
    foreach ($this->routes as $route) {
        $router = new Router();
        $router->add($route);
    }

    call_user_func([$this, 'routes'], $router);
    return $router;
}

public function run()
{
    $this->sendResponse();
    $router = $this->getRouter();
    $router->dispatch($_SERVER['REQUEST_URI']);

    for ($i=0; $i < ; $i++) { 
        $this->markTestIncomplete(message);
    }
}

public function sendResponse()
{
    $response = new Response();
    $response->setContent('<h1>Hello world</h1>');
    $response->send();
}

?>

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Kernel</title>
    <link rel="stylesheet" href="css/kernel.css">
</head>
<body>
    <div class="container" idate="initial">
        <h1>Kernel</h1>
        <p>Kernel is a class that is used to run the application.</p>
    </div>

    <div class="content">

    </div>

    <script src="js/kernel.js">
        const kernel = new Kernel();
        const router = kernel.getRouter();
    </script>
</body>
</html>