module.exports = {
  'env': {
    'browser': true,
    'es2021': true
  },
  'extends': [
    'eslint:recommended',
    'plugin:vue/essential',
    'plugin:@typescript-eslint/recommended'
  ],
  // parser: 指定解析器
  // yarn add -D esprima -- ESLint 默认安装并使用
  // yarn add -D @babel/eslint-parser
  // yarn add -D @typescript-eslint/parser
  // yarn add -D eslint-plugin-vue -- 取其依赖 vue-eslint-parser
  'parser': 'vue-eslint-parser',
  // parserOptions: JavaScript 语言选项
  'parserOptions': {
    'ecmaVersion': 'latest',
    'parser': '@typescript-eslint/parser',
    'sourceType': 'module'
  },
  'plugins': [
    'vue',
    '@typescript-eslint'
  ],
  'globals': {
    defineProps: 'readonly',
    defineEmits: 'readonly',
    defineExpose: 'readonly',
    withDefaults: 'readonly',
  },
  'rules': {
    'semi': ['warn', 'never'],
    'quotes': ['error', 'single'],
    'vue/no-v-model-argument': 'off',
  },
}
