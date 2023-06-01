## Set up

---

Instalar dependencias antes de correr el proyecto:

```bash
npm i
```

Para correr el proyecto:

```bash
npm run dev
```

Para ver datos de SQLite:

```bash
npx prisma studio
```

Para hacer una nueva migration:

```bash
npm prisma migrate dev
```

- `dev` hace referencia a el archivo `dev.db` que es la base de datos SQLite.

## Endpoints

---

### /api/measures

#### POST

##### Payload

```
{
	data: String;
	filtration: Boolean;
}
```

https://remaster.com/blog/next-auth-jwt-session
