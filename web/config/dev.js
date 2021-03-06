import webpack from 'webpack';
import merge from 'webpack-merge';
import baseConfig from './base';

export default merge(baseConfig, {
  devtool: 'cheap-module-eval-source-map',

  output: {
    filename: '[name].js',
    publicPath: '/assets/'
  },

  entry: {
    home: ['react-hot-loader', 'babel-polyfill', './app/Home/index.jsx'],
    register: ['react-hot-loader', 'babel-polyfill', './app/Register/index.jsx'],
    login: ['react-hot-loader', 'babel-polyfill', './app/Login/index.jsx']
  },

  module: {
    loaders: [
      {
        test: /\.jsx?$/,
        exclude: /node_modules/,
        loaders: [
          'babel-loader',
          'eslint-loader'
        ]
      },
      {
        test: /\.s?css$/,
        loaders: [
          'style-loader',
          'css-loader?!postcss-loader!sass-loader'
        ]
      },
      {
        test: /\.(jpg|png|gif|svg)$/,
        loader: 'file-loader'
      },
      {
        test: /\.(woff|woff2|eot|ttf)$/,
        loader: 'url-loader?limit=100000'
      }
    ]
  },

  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin(),
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify('development')
    }),
  ],

  devServer: {
    hot: true,
    publicPath: '/assets/',
    historyApiFallback: true,
    proxy: {
      '/!(assets)**': 'http://localhost:4000',
    },
    host: '0.0.0.0'
  }
});
