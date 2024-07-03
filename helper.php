<?php

use App\Model\Vod;
use App\Model\VodType;
use Hyperf\Database\Model\Collection;

if (!function_exists("queryVods")) {
    function queryVods(array $parameters)
    {
        if (empty($parameters)) {
            return Collection::empty();
        }
        $query = Vod::query()->select($parameters["fields"] ?? ["*"]);
        foreach ($parameters["where"] ?? [] as $k => $v) {
            if ($v[0] == "in") {
                $boolean = !empty($v[2]['or']) ? 'or' : 'and';
                $not = !empty($v[2]['not']) ? true : false;
                $query->whereIn($k, $v[1], $boolean, $not);
                continue;
            }
            if ($v[0] == "between") {
                $boolean = !empty($v[2]['or']) ? 'or' : 'and';
                $not = !empty($v[2]['not']) ? true : false;
                $query->whereBetween($v, $v[1], $boolean, $not);
                continue;
            }
            if ($v[0] === null) {
                $boolean = !empty($v[1]['or']) ? 'or' : 'and';
                $not = !empty($v[1]['not']) ? true : false;
                $query->whereNull($k, $boolean, $not);
                continue;
            }

            $boolean = !empty($v[1]['or']) ? 'or' : 'and';
            $not = !empty($v[1]['not']) ? true : false;
            $query->where($k, $v[0], $v[1], $boolean, $not);
        }
        foreach ($parameters['order'] ?? [] as   $order) {
            $query->orderBy($order[0], $order[1]);
        
        $limit = $parameters['limit'] ?? 10;
        $page = $parameters['page'] ?? 1;
        $vods =  $query->forPage($page, $limit)->get();
        return   $vods;
    }
}
