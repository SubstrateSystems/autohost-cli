```txt
    
    / \  _   _| |_ ___| |__   ___  ___| |_   / __| |   |_ _|
   / _ \| | | | __/ _ \ '_ \ / _ \/ __| __| | |  | |    | | 
  / ___ \ |_| | ||  _ \ | | | (_) \__ \ |_  | |__| |___ | | 
 /_/   \_\__,_|\__\___|_| |_|\___/|___/\__|  \___|_____|___|                        
```

# 🚀 AutoHost CLI

**Recupera el control de tus servicios.**  
**AutoHost CLI** es una herramienta de línea de comandos para instalar, configurar y administrar aplicaciones y servicios **en tu propio servidor**, sin depender de terceros y con un flujo de trabajo sencillo y automatizado.

---

## 🌟 Características

- **Instalación en un comando**: Despliega aplicaciones listas para usar con `app install`.
- **Soporte para múltiples apps**: Nextcloud, BookStack, Redis, MySQL y más (¡en constante crecimiento!).
- **Integración con Tailscale**: Conéctate de forma segura a tu infraestructura privada.
- **Compatibilidad con Docker**: Aislamiento y portabilidad de tus aplicaciones.
- **Enfoque en privacidad y control**: Todo se ejecuta en **tu** infraestructura.

---

## ⚙️ Requisitos Previos

Antes de instalar, asegúrate de contar con:
- Un sistema basado en **Linux** (compatible con distribuciones modernas como Ubuntu/Debian).  
- **Docker** instalado y corriendo.  
- Permisos de administrador (**sudo/root**).  
- Opcional: cuenta de **Tailscale** si quieres habilitar acceso seguro privado.  

---

## 📦 Instalación

Instala AutoHost CLI directamente desde GitHub con un solo comando:

```bash
curl -fsSL https://raw.githubusercontent.com/mazapanuwu13/autohost-cli/main/scripts/install.sh | bash
```

Este script detecta automáticamente tu sistema operativo y arquitectura, descarga la versión más reciente del binario desde GitHub Releases e instala AutoHost CLI en tu sistema.

---

## 🛠 Uso Básico

### Flujo de ejemplo

```bash
# Inicializar entorno
autohost init

# Configuración inicial (dominio, redes, etc.)
autohost setup

# Instalar una aplicación (ejemplo: Nextcloud)
autohost app install

# Levantar la aplicación
autohost app start nextcloud

# Ver estado de la app
autohost app status nextcloud
```

---

## 📂 Aplicaciones soportadas

| App        | Puerto por defecto | Estado  |
|------------|-------------------|---------|
| Nextcloud  | 8081              | ✅ Estable |
| BookStack  | 6875              | ✅ Estable |
| MySQL      | 3306              | ✅ Estable |

*(La lista crece con cada versión. ¡Tu feedback ayuda a priorizar nuevas apps!)*

---

## 🔒 Filosofía

En un mundo donde la mayoría de las aplicaciones están en la nube, **AutoHost CLI** te devuelve el poder:  
- Controlas **tus datos**.  
- Eliminas la dependencia de múltiples SaaS.  
- Construyes tu propia infraestructura, escalable y privada.  

---

## 🧪 Desarrollo y Pruebas

### Ejecutar pruebas

Para correr todas las pruebas del proyecto:

```bash
go test ./...
```

Para ejecutar pruebas con información detallada:

```bash
go test -v ./...
```

Para ejecutar pruebas de un paquete específico:

```bash
go test -v ./utils
go test -v ./internal/adapters/cli/app
```

### Estructura de pruebas

- Las pruebas siguen la convención `*_test.go` junto a los archivos de código
- Se incluyen pruebas unitarias para los comandos principales (`app` y sus subcomandos)
- Las pruebas utilizan el paquete estándar `testing` de Go
- Los mocks se utilizan para componentes externos (Docker, sistema de archivos)

### Cobertura de pruebas

Las pruebas actuales cubren:
- ✅ Comandos del CLI (estructura y validación de argumentos)
- ✅ Utilidades de helpers (`ValidPort`, `GetAutohostDir`, etc.)
- ✅ Validación de configuraciones básicas

---

## 🤝 Contribuir

¿Quieres aportar?  
1. Haz un fork del repositorio.  
2. Crea una rama para tu feature/fix.  
3. Envía un Pull Request.  
4. Revisa las issues con la etiqueta **good first issue** para comenzar.

---

## 📜 Licencia

Este proyecto está bajo la licencia **MIT**.

---

> 💡 **Consejo:** Si quieres recibir actualizaciones y novedades, visita [autohst.dev](https://autohst.dev) o síguenos en redes.