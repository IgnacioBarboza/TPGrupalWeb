// ============================================================================
// 1. FILM (/films)
// ============================================================================
async function listarFilms() {
  const contenedor = document.getElementById("list-films");

  try {
    const respuesta = await fetch("/films");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const films = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(films) && films.length > 0) {
      const fragmento = document.createDocumentFragment();

      films.forEach((film) => {
        const li = document.createElement("li");
        const fecha = film.filmreleasedate
          ? film.filmreleasedate.split("T")[0]
          : "";

        const texto = document.createElement("span");
        texto.textContent = `[ID: ${film.filmid}] ${film.filmtitle} - Estreno: ${fecha} - Duración: ${film.filmduration} min `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () => eliminarFilm(film.filmid));

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML = "<li>No hay películas cargadas.</li>";
    }
  } catch (error) {
    console.error("Error al obtener las películas:", error);
    contenedor.innerHTML = `<li>Error al cargar películas: ${error.message}</li>`;
  }
}

async function crearFilm(event) {
  event.preventDefault();

  const titulo = document.getElementById("film-title").value.trim();
  const fechaInput = document.getElementById("film-release-date").value;
  const duracion = parseInt(document.getElementById("film-duration").value, 10);

  const nuevoFilm = {
    filmtitle: titulo,
    filmreleasedate: `${fechaInput}T00:00:00Z`,
    filmduration: duracion,
  };

  try {
    const respuesta = await fetch("/films/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevoFilm),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-film").reset();
    await listarFilms();
  } catch (error) {
    console.error("Error al crear la película:", error);
    alert(`No se pudo crear la película: ${error.message}`);
  }
}

async function eliminarFilm(id) {
  try {
    const respuesta = await fetch(`/films/${id}`, { method: "DELETE" });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarFilms();
  } catch (error) {
    console.error(`Error al eliminar la película ${id}:`, error);
    alert(`No se pudo eliminar la película: ${error.message}`);
  }
}

// ============================================================================
// 2. GENRE (/genres)
// ============================================================================
async function listarGenres() {
  const contenedor = document.getElementById("list-genres");

  try {
    const respuesta = await fetch("/genres");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const genres = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(genres) && genres.length > 0) {
      const fragmento = document.createDocumentFragment();

      genres.forEach((genre) => {
        const li = document.createElement("li");

        const texto = document.createElement("span");
        texto.textContent = `[ID: ${genre.genreid}] ${genre.genrename} `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () =>
          eliminarGenre(genre.genreid)
        );

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML = "<li>No hay géneros cargados.</li>";
    }
  } catch (error) {
    console.error("Error al obtener los géneros:", error);
    contenedor.innerHTML = `<li>Error al cargar géneros: ${error.message}</li>`;
  }
}

async function crearGenre(event) {
  event.preventDefault();

  const nombre = document.getElementById("genre-name").value.trim();

  try {
    const respuesta = await fetch("/genres/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ genrename: nombre }),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-genre").reset();
    await listarGenres();
  } catch (error) {
    console.error("Error al crear el género:", error);
    alert(`No se pudo crear el género: ${error.message}`);
  }
}

async function eliminarGenre(id) {
  try {
    const respuesta = await fetch(`/genres/${id}`, { method: "DELETE" });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarGenres();
  } catch (error) {
    console.error(`Error al eliminar el género ${id}:`, error);
    alert(`No se pudo eliminar el género: ${error.message}`);
  }
}

// ============================================================================
// 3. CAST MEMBER (/castmembers)
// ============================================================================
async function listarCastMembers() {
  const contenedor = document.getElementById("list-castmembers");

  try {
    const respuesta = await fetch("/castmembers");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const castMembers = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(castMembers) && castMembers.length > 0) {
      const fragmento = document.createDocumentFragment();

      castMembers.forEach((cm) => {
        const li = document.createElement("li");

        const texto = document.createElement("span");
        texto.textContent = `[ID: ${cm.castmemberid}] ${cm.castmembername} - Rol: ${cm.castmemberrole} - Género: ${cm.castmembergender} `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () =>
          eliminarCastMember(cm.castmemberid)
        );

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML = "<li>No hay miembros del elenco cargados.</li>";
    }
  } catch (error) {
    console.error("Error al obtener miembros del elenco:", error);
    contenedor.innerHTML = `<li>Error al cargar miembros: ${error.message}</li>`;
  }
}

async function crearCastMember(event) {
  event.preventDefault();

  const nuevoMiembro = {
    castmembername: document.getElementById("castmember-name").value.trim(),
    castmemberrole: document.getElementById("castmember-role").value.trim(),
    castmembergender: document.getElementById("castmember-gender").value.trim(),
  };

  try {
    const respuesta = await fetch("/castmembers/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevoMiembro),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-castmember").reset();
    await listarCastMembers();
  } catch (error) {
    console.error("Error al crear miembro del elenco:", error);
    alert(`No se pudo crear el miembro: ${error.message}`);
  }
}

async function eliminarCastMember(id) {
  try {
    const respuesta = await fetch(`/castmembers/${id}`, { method: "DELETE" });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarCastMembers();
  } catch (error) {
    console.error(`Error al eliminar miembro ${id}:`, error);
    alert(`No se pudo eliminar el miembro: ${error.message}`);
  }
}

// ============================================================================
// 4. STATUS (/statuses)
// ============================================================================
async function listarStatuses() {
  const contenedor = document.getElementById("list-statuses");

  try {
    const respuesta = await fetch("/statuses");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const statuses = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(statuses) && statuses.length > 0) {
      const fragmento = document.createDocumentFragment();

      statuses.forEach((st) => {
        const li = document.createElement("li");

        const texto = document.createElement("span");
        texto.textContent = `[ID: ${st.statusid}] ${st.statusname} `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () =>
          eliminarStatus(st.statusid)
        );

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML = "<li>No hay estados cargados.</li>";
    }
  } catch (error) {
    console.error("Error al obtener estados:", error);
    contenedor.innerHTML = `<li>Error al cargar estados: ${error.message}</li>`;
  }
}

async function crearStatus(event) {
  event.preventDefault();

  const nombre = document.getElementById("status-name").value.trim();

  try {
    const respuesta = await fetch("/statuses/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ statusname: nombre }),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-status").reset();
    await listarStatuses();
  } catch (error) {
    console.error("Error al crear estado:", error);
    alert(`No se pudo crear el estado: ${error.message}`);
  }
}

async function eliminarStatus(id) {
  try {
    const respuesta = await fetch(`/statuses/${id}`, { method: "DELETE" });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarStatuses();
  } catch (error) {
    console.error(`Error al eliminar estado ${id}:`, error);
    alert(`No se pudo eliminar el estado: ${error.message}`);
  }
}

// ============================================================================
// 5. USER (/users)
// ============================================================================
async function listarUsers() {
  const contenedor = document.getElementById("list-users");

  try {
    const respuesta = await fetch("/users");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const users = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(users) && users.length > 0) {
      const fragmento = document.createDocumentFragment();

      users.forEach((user) => {
        const li = document.createElement("li");
        const cumple = user.userbirthday
          ? user.userbirthday.split("T")[0]
          : "";

        const texto = document.createElement("span");
        texto.textContent = `[ID: ${user.userid}] ${user.username} - Nacimiento: ${cumple} - Género: ${user.usergenre} `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () => eliminarUser(user.userid));

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML = "<li>No hay usuarios cargados.</li>";
    }
  } catch (error) {
    console.error("Error al obtener usuarios:", error);
    contenedor.innerHTML = `<li>Error al cargar usuarios: ${error.message}</li>`;
  }
}

async function crearUser(event) {
  event.preventDefault();

  const fechaInput = document.getElementById("user-birthday").value;

  const nuevoUser = {
    username: document.getElementById("user-name").value.trim(),
    userpasswordhashed: document.getElementById("user-password").value,
    userbirthday: `${fechaInput}T00:00:00Z`,
    usergenre: document.getElementById("user-genre").value.trim(),
  };

  try {
    const respuesta = await fetch("/users/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevoUser),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-user").reset();
    await listarUsers();
  } catch (error) {
    console.error("Error al crear usuario:", error);
    alert(`No se pudo crear el usuario: ${error.message}`);
  }
}

async function eliminarUser(id) {
  try {
    const respuesta = await fetch(`/users/${id}`, { method: "DELETE" });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarUsers();
  } catch (error) {
    console.error(`Error al eliminar usuario ${id}:`, error);
    alert(`No se pudo eliminar el usuario: ${error.message}`);
  }
}

// ============================================================================
// 6. REVIEW (/reviews) - Clave compuesta: (reviewfilmid, reviewuserid)
// ============================================================================
async function listarReviews() {
  const contenedor = document.getElementById("list-reviews");

  try {
    const respuesta = await fetch("/reviews");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const reviews = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(reviews) && reviews.length > 0) {
      const fragmento = document.createDocumentFragment();

      reviews.forEach((rev) => {
        const li = document.createElement("li");

        const texto = document.createElement("span");
        texto.textContent = `[Película ID: ${rev.reviewfilmid} | Usuario ID: ${rev.reviewuserid}] ${rev.reviewstars}★ - "${rev.reviewdescription}" `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () =>
          eliminarReview(rev.reviewfilmid, rev.reviewuserid)
        );

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML = "<li>No hay reseñas cargadas.</li>";
    }
  } catch (error) {
    console.error("Error al obtener reseñas:", error);
    contenedor.innerHTML = `<li>Error al cargar reseñas: ${error.message}</li>`;
  }
}

async function crearReview(event) {
  event.preventDefault();

  const nuevaReview = {
    reviewfilmid: parseInt(document.getElementById("review-film-id").value, 10),
    reviewuserid: parseInt(document.getElementById("review-user-id").value, 10),
    reviewstars: parseInt(document.getElementById("review-stars").value, 10),
    reviewdescription: document
      .getElementById("review-description")
      .value.trim(),
  };

  try {
    const respuesta = await fetch("/reviews/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevaReview),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-review").reset();
    await listarReviews();
  } catch (error) {
    console.error("Error al crear reseña:", error);
    alert(`No se pudo crear la reseña: ${error.message}`);
  }
}

async function eliminarReview(filmId, userId) {
  try {
    const respuesta = await fetch(`/reviews/${filmId}/${userId}`, {
      method: "DELETE",
    });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarReviews();
  } catch (error) {
    console.error(`Error al eliminar reseña (${filmId}, ${userId}):`, error);
    alert(`No se pudo eliminar la reseña: ${error.message}`);
  }
}

// ============================================================================
// 7. FILM CAST (/filmcasts) - Clave compuesta: (filmid, castmemberid)
// ============================================================================
async function listarFilmCasts() {
  const contenedor = document.getElementById("list-filmcasts");

  try {
    const respuesta = await fetch("/filmcasts");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const filmCasts = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(filmCasts) && filmCasts.length > 0) {
      const fragmento = document.createDocumentFragment();

      filmCasts.forEach((fc) => {
        const li = document.createElement("li");

        const texto = document.createElement("span");
        texto.textContent = `Película ID: ${fc.filmid} ↔ Miembro Elenco ID: ${fc.castmemberid} `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () =>
          eliminarFilmCast(fc.filmid, fc.castmemberid)
        );

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML = "<li>No hay relaciones FilmCast cargadas.</li>";
    }
  } catch (error) {
    console.error("Error al obtener FilmCasts:", error);
    contenedor.innerHTML = `<li>Error al cargar FilmCasts: ${error.message}</li>`;
  }
}

async function crearFilmCast(event) {
  event.preventDefault();

  const nuevoFilmCast = {
    filmid: parseInt(document.getElementById("filmcast-film-id").value, 10),
    castmemberid: parseInt(
      document.getElementById("filmcast-castmember-id").value,
      10
    ),
  };

  try {
    const respuesta = await fetch("/filmcasts/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevoFilmCast),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-filmcast").reset();
    await listarFilmCasts();
  } catch (error) {
    console.error("Error al crear FilmCast:", error);
    alert(`No se pudo crear la relación FilmCast: ${error.message}`);
  }
}

async function eliminarFilmCast(filmId, castMemberId) {
  try {
    const respuesta = await fetch(`/filmcasts/${filmId}/${castMemberId}`, {
      method: "DELETE",
    });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarFilmCasts();
  } catch (error) {
    console.error(
      `Error al eliminar FilmCast (${filmId}, ${castMemberId}):`,
      error
    );
    alert(`No se pudo eliminar la relación FilmCast: ${error.message}`);
  }
}

// ============================================================================
// 8. FILM GENRE (/filmgenres) - Clave compuesta: (genre_genreid, film_filmid)
// ============================================================================
async function listarFilmGenres() {
  const contenedor = document.getElementById("list-filmgenres");

  try {
    const respuesta = await fetch("/filmgenres");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const filmGenres = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(filmGenres) && filmGenres.length > 0) {
      const fragmento = document.createDocumentFragment();

      filmGenres.forEach((fg) => {
        const li = document.createElement("li");

        const texto = document.createElement("span");
        texto.textContent = `Género ID: ${fg.genre_genreid} ↔ Película ID: ${fg.film_filmid} `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () =>
          eliminarFilmGenre(fg.genre_genreid, fg.film_filmid)
        );

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML = "<li>No hay relaciones FilmGenre cargadas.</li>";
    }
  } catch (error) {
    console.error("Error al obtener FilmGenres:", error);
    contenedor.innerHTML = `<li>Error al cargar FilmGenres: ${error.message}</li>`;
  }
}

async function crearFilmGenre(event) {
  event.preventDefault();

  const nuevoFilmGenre = {
    genre_genreid: parseInt(
      document.getElementById("filmgenre-genre-id").value,
      10
    ),
    film_filmid: parseInt(
      document.getElementById("filmgenre-film-id").value,
      10
    ),
  };

  try {
    const respuesta = await fetch("/filmgenres/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevoFilmGenre),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-filmgenre").reset();
    await listarFilmGenres();
  } catch (error) {
    console.error("Error al crear FilmGenre:", error);
    alert(`No se pudo crear la relación FilmGenre: ${error.message}`);
  }
}

async function eliminarFilmGenre(genreId, filmId) {
  try {
    const respuesta = await fetch(`/filmgenres/${genreId}/${filmId}`, {
      method: "DELETE",
    });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarFilmGenres();
  } catch (error) {
    console.error(`Error al eliminar FilmGenre (${genreId}, ${filmId}):`, error);
    alert(`No se pudo eliminar la relación FilmGenre: ${error.message}`);
  }
}

// ============================================================================
// 9. FILM USER STATUS (/filmuserstatuses) - Clave compuesta: (film_filmid, user_userid)
// ============================================================================
async function listarFilmUserStatuses() {
  const contenedor = document.getElementById("list-filmuserstatuses");

  try {
    const respuesta = await fetch("/filmuserstatuses");
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }

    const statuses = await respuesta.json();
    contenedor.innerHTML = "";

    if (Array.isArray(statuses) && statuses.length > 0) {
      const fragmento = document.createDocumentFragment();

      statuses.forEach((fus) => {
        const li = document.createElement("li");

        const texto = document.createElement("span");
        texto.textContent = `Película ID: ${fus.film_filmid} | Usuario ID: ${fus.user_userid} → Estado ID: ${fus.status_statusid} `;

        const btnEliminar = document.createElement("button");
        btnEliminar.type = "button";
        btnEliminar.textContent = "Eliminar";
        btnEliminar.addEventListener("click", () =>
          eliminarFilmUserStatus(fus.film_filmid, fus.user_userid)
        );

        li.appendChild(texto);
        li.appendChild(btnEliminar);
        fragmento.appendChild(li);
      });

      contenedor.appendChild(fragmento);
    } else {
      contenedor.innerHTML =
        "<li>No hay estados de película por usuario cargados.</li>";
    }
  } catch (error) {
    console.error("Error al obtener FilmUserStatuses:", error);
    contenedor.innerHTML = `<li>Error al cargar FilmUserStatuses: ${error.message}</li>`;
  }
}

async function crearFilmUserStatus(event) {
  event.preventDefault();

  const nuevoFilmUserStatus = {
    film_filmid: parseInt(document.getElementById("fus-film-id").value, 10),
    user_userid: parseInt(document.getElementById("fus-user-id").value, 10),
    status_statusid: parseInt(
      document.getElementById("fus-status-id").value,
      10
    ),
  };

  try {
    const respuesta = await fetch("/filmuserstatuses/", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(nuevoFilmUserStatus),
    });

    if (!respuesta.ok) {
      const detalle = await respuesta.text();
      throw new Error(detalle || `Error HTTP: ${respuesta.status}`);
    }

    document.getElementById("form-filmuserstatus").reset();
    await listarFilmUserStatuses();
  } catch (error) {
    console.error("Error al crear FilmUserStatus:", error);
    alert(`No se pudo asignar el estado: ${error.message}`);
  }
}

async function eliminarFilmUserStatus(filmId, userId) {
  try {
    const respuesta = await fetch(`/filmuserstatuses/${filmId}/${userId}`, {
      method: "DELETE",
    });
    if (!respuesta.ok) {
      throw new Error(`Error HTTP: ${respuesta.status}`);
    }
    await listarFilmUserStatuses();
  } catch (error) {
    console.error(
      `Error al eliminar FilmUserStatus (${filmId}, ${userId}):`,
      error
    );
    alert(`No se pudo eliminar el estado asignado: ${error.message}`);
  }
}

// ============================================================================
// INICIALIZACIÓN AL CARGAR EL DOM
// ============================================================================
document.addEventListener("DOMContentLoaded", () => {
  // Cargar todas las listas iniciales
  listarFilms();
  listarGenres();
  listarCastMembers();
  listarStatuses();
  listarUsers();
  listarReviews();
  listarFilmCasts();
  listarFilmGenres();
  listarFilmUserStatuses();

  // Asociar eventos submit a cada formulario
  document.getElementById("form-film")?.addEventListener("submit", crearFilm);
  document.getElementById("form-genre")?.addEventListener("submit", crearGenre);
  document
    .getElementById("form-castmember")
    ?.addEventListener("submit", crearCastMember);
  document
    .getElementById("form-status")
    ?.addEventListener("submit", crearStatus);
  document.getElementById("form-user")?.addEventListener("submit", crearUser);
  document
    .getElementById("form-review")
    ?.addEventListener("submit", crearReview);
  document
    .getElementById("form-filmcast")
    ?.addEventListener("submit", crearFilmCast);
  document
    .getElementById("form-filmgenre")
    ?.addEventListener("submit", crearFilmGenre);
  document
    .getElementById("form-filmuserstatus")
    ?.addEventListener("submit", crearFilmUserStatus);
});
