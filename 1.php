<?php
//https://pdai.tech/md/db/nosql-es/elasticsearch-x-index-mapping.html
$servers = [
    [
        "url" => "http://127.0.0.1:8898/admin/reverseproxy",//后台地址
        "username" => "",//后台用户名
        "password" => "",//后台密码
        "private_js"=>"",//服务器私有的js代码，例如百度统计代码
    ],
];
$jsContent = file_get_contents("inject.js");

foreach ($servers as $server) {
    $jsContent=$server['private_js']."\n\n".$jsContent;
    $param=[
        'username'=>$server['username'],
        'password'=>$server['password'],
        'js_content'=>base64_encode($jsContent);
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
