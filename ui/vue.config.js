module.exports = {
  productionSourceMap: false,
  configureWebpack: {
    devtool: "source-map"
  },
  outputDir: "assets",
  transpileDependencies: ["vuetify"]
};
