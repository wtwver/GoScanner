<?php
    $_POST = json_decode(file_get_contents('php://input'), true);

    $items = array(
        "activeCrawler" => $_POST["activeCrawler"],
        "activeFuzzer" => $_POST["activeFuzzer"],
        "passiveCrawler" => $_POST["passiveCrawler"],
        "urlSeed" => $_POST["urlSeed"]
    );
    
    header("Content-type: application/json");
    echo json_encode($items);
?>