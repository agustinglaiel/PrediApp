// .eslintrc.js
module.exports = {
  env: {
    browser: true,
    es2021: true,
  },
  extends: [
    "eslint:recommended",
    "plugin:react/recommended",
    "plugin:tailwindcss/recommended", // Agrega el plugin de Tailwind
  ],
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 12,
    sourceType: "module",
  },
  plugins: [
    "react",
    "tailwindcss", // Asegúrate de incluir el plugin
  ],
  rules: {
    "tailwindcss/classnames-order": "warn", // Ordena las clases de Tailwind (opcional)
    "react/react-in-jsx-scope": "off",
  },
  settings: {
    react: {
      version: "detect", // Detecta automáticamente la versión de React
    },
  },
};
