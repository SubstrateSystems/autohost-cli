# Database Seeding

Este directorio contiene los archivos de seeding para poblar la base de datos con datos iniciales.

## Estructura

- Los archivos de seeding deben seguir el patrón `XXX_nombre.sql` donde `XXX` es un número de orden (ej: `001_catalog_apps.sql`)
- Los seeds se ejecutan automáticamente después de las migraciones al inicializar la aplicación
- El sistema trackea qué seeds ya se han ejecutado usando la tabla `_seeds`
- Los seeds son idempotentes - pueden ejecutarse múltiples veces sin causar errores

## Seeds Actuales

### 001_catalog_apps.sql
Inserta las aplicaciones disponibles en el catálogo:
- bookstack: Plataforma para organizar y almacenar información
- nextcloud: Suite de software para servicios de hosting de archivos  
- redis: Base de datos en memoria para cache y mensajería

## Agregar Nuevos Seeds

1. Crear un archivo con formato `XXX_nombre.sql` en este directorio
2. Usar `INSERT OR IGNORE` para evitar duplicados
3. Los seeds se ejecutarán automáticamente en el próximo arranque de la aplicación

## Notas

- Los seeds se ejecutan en orden alfabético de nombre de archivo
- Una vez ejecutado, un seed no se vuelve a ejecutar
- Si necesitas modificar datos existentes, crea un nuevo seed con un número mayor
