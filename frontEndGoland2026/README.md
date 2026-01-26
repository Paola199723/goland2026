# Frontend Goland 2026

Frontend moderno desarrollado con Vue 3, TypeScript, Pinia y Tailwind CSS.

## Características

- ✅ Autenticación de usuarios
- ✅ Visualización de challenges en tabla
- ✅ Paginación
- ✅ Integración con API backend
- ✅ Diseño responsivo con Tailwind CSS
- ✅ Gestión de estado con Pinia
- ✅ TypeScript para type safety

## Requisitos

- Node.js 16+
- npm o yarn

## Instalación

```bash
cd frontEndGoland2026
npm install
```

## Desarrollo

```bash
npm run dev
```

La aplicación estará disponible en `http://localhost:5173`

## Build para producción

```bash
npm run build
```

## Estructura del Proyecto

```
src/
├── pages/              # Páginas principales
│   ├── LoginPage.vue   # Página de login
│   └── DashboardPage.vue # Panel principal
├── components/         # Componentes reutilizables
├── services/          # Servicios API
├── stores/            # Store de Pinia
├── App.vue            # Componente raíz
├── main.ts            # Punto de entrada
└── style.css          # Estilos globales
```

## Variables de Entorno

Configurar en un archivo `.env` si es necesario:

```
VITE_API_BASE_URL=http://localhost:8080
```
