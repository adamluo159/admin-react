var webpack = require('webpack');
var path = require('path');

var commonsPlugin = new webpack.optimize.CommonsChunkPlugin('common.js');
var hot = new webpack.HotModuleReplacementPlugin();
 
module.exports = {
    devServer:{
        colors : true,
        historyApiFallback: true,
        inline: true
    },
    //devtool : 'source-map',
    //插件项
    plugins: [commonsPlugin, hot],
    //页面入口文件配置
    entry: [
        './src/index.js',
        "webpack-dev-server/client?http://0.0.0.0:3000",
        "webpack/hot/only-dev-server"
        //'webpack-dev-server/client?http://127.0.0.1:3000',
        //'webpack/hot/only-dev-server'
    ],
    //入口文件输出配置
    output: {
        path: __dirname + '/dist',
        filename: 'bundle.js',
        publicPath: '/src/'
    },
    module: {
        //加载器配置
        loaders: [
            { test: /\.css$/, loader: 'style-loader!css-loader' },
            //{ test: /\.js$/, loaders: ['babel'], exclude: /node_modules/},
            {test: /\.js$/, loaders: ['react-hot', 'babel?presets[]=es2015,presets[]=react,presets[]=stage-0'],exclude: /node_modules/},
            { test: /\.less$/, loader: 'style!css!less'},
            { test: /\.(png|jpg)$/, loader: 'url-loader?limit=8192'}
        ]
    },
    //其它解决方案配置
    resolve: {
        //root: 'E:/github/flux-example/src', //绝对路径
        extensions: ['', '.js', '.json', '.scss'],
        //alias: {
        //    AppStore : 'js/stores/AppStores.js',
        //    ActionType : 'js/actions/ActionType.js',
        //    AppAction : 'js/actions/AppAction.js'
        //}
    }
};
