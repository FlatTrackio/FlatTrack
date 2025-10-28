module.exports = {
  module: {
    rules: [
      {
        test: /.scss$/,
        use: [
          "style-loader",
          "css-loader",
          {
            loader: "sass-loader",
            options: {
              implementation: require("sass"), // Dart Sass
              sassOptions: {
                fiber: false, // Removed in Dart Sass 2.0
              },
            },
          },
        ],
      },
    ],
  },
};
