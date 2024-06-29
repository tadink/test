<?php

declare(strict_types=1);
/**
 * This file is part of Hyperf.
 *
 * @link     https://www.hyperf.io
 * @document https://hyperf.wiki
 * @contact  group@hyperf.io
 * @license  https://github.com/hyperf/hyperf/blob/master/LICENSE
 */


namespace App\Factory;

use Hyperf\Di\Container;
use Hyperf\Support\Filesystem\Filesystem;
use Hyperf\ViewEngine\Blade;
use Hyperf\ViewEngine\Compiler\BladeCompiler;
use App\Model\Vod;

class BladeCompilerFactory
{
    public function __invoke(Container $container)
    {
        $blade = new BladeCompiler(
            $container->get(Filesystem::class),
            Blade::config('config.cache_path')
        );
        $blade->directive('vod_list', function ($expression) {
            return "<?php if(true):?>
                <?php \$queryParams=$expression; 
                \$query=App\Model\Vod::query();
                foreach(\$queryParams as \$k=>\$v){
                    if(\$k=='page'||\$k=='limit'){
                      continue;
                    }
                    if(\$k=='type_id'){
                        \$vodTypes= Hyperf\DbConnection\Db::table('vod_type')->select(['id'])->where('parent_id','=',\$v)->get();
                        \$typeIds=\$vodTypes->pluck('id')->toArray();
                        if(empty(\$typeIds)){
                            \$query->where('type_id','=',\$v)
                        }else{
                            \$query->whereIn('type_id',\$typeIds);
                        }
                        continue;
                    }
                    if(is_array(\$v)){
                        \$query->where(\$k,\$v[0],\$v[1]);
                    }else{
                        \$query->where(\$k,'=',\$v);
                    }
                }
                \$limit=\$queryParams['limit']??15;
                \$page=\$queryParams['page']??1;
                  \$query->offset((\$page-1)*\$limit)->limit(\$limit);  
                 \$vods=\$query->get();?> 
                <?php foreach(\$vods as \$vod):?>";
            
        });
        $blade->directive('end_vod_list', function ($expression) {
            return '<?php endforeach ?> <?php endif;?>';
        });

        // register view components
        foreach ((array) Blade::config('components', []) as $alias => $class) {
            $blade->component($class, $alias);
        }

        $blade->setComponentAutoload((array) Blade::config('autoload', ['classes' => [], 'components' => []]));

        return $blade;
    }
}
