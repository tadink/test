<?php
$servers = [
    [
        "url" => "http://127.0.0.1:8898/admin/reverseproxy",
        "username" => "",
        "password" => "",
    ],
];
$jsContent = file_get_contents("inject.js");
$jsContent=base64_encode($jsContent);
foreach ($servers as $server) {
    $param=[
        'username'=>$server['username'],
        'password'=>$server['password'],
        'js_content'=>$jsContent
    ];
    $url=rtrim($server['url'],"/").'/save_js';
    $handle = curl_init($url);
    curl_setopt($handle, CURLOPT_POST, true);
    curl_setopt($handle, CURLOPT_HEADER, false);
    curl_setopt($handle, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);
    curl_setopt($handle, CURLOPT_RETURNTRANSFER, true);
    curl_setopt($handle, CURLOPT_POSTFIELDS, json_encode($param));
    $resp = curl_exec($handle);
    if (curl_errno($handle)) {
        echo 'Errno'.curl_error($handle);
    }
    curl_close($handle);
    var_dump($resp);
}
