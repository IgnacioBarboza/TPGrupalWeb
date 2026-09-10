-- name: CreateFilm :exec
INSERT INTO Film(FilmId,FilmTitle,FilmReleaseDate,FilmDuration) 
VALUES($1,$2,$3,$4)
RETURNING FilmId,FilmTitle,FilmReleaseDate,FilmDuration;

-- name: GetFilm :one
SELECT FilmId,FilmTitle,FilmReleaseDate,FilmDuration
FROM Film
WHERE FilmId = $1;

-- name: GetFilms :many
SELECT FilmId,FilmTitle,FilmReleaseDate,FilmDuration
FROM Film;

-- name: UpdateFilm :exec
UPDATE Film
SET FilmTitle = $2 ,FilmReleaseDate = $3, FilmDuration = $4
WHERE FilmId = $1;

-- name: DeleteFilm :exec
DELETE FROM Film
WHERE FilmId = $1;

-- name: CreateCastMember :exec
INSERT INTO CastMember(CastMemberId, CastMemberName, CastMemberRole, CastMemberGender)
VALUES($1, $2, $3, $4)
RETURNING CastMemberId, CastMemberName, CastMemberRole, CastMemberGender;

-- name: GetCastMember :one
SELECT CastMemberId, CastMemberName, CastMemberRole, CastMemberGender
FROM CastMember
WHERE CastMemberId = $1;

-- name: GetCastMembers :many
SELECT CastMemberId, CastMemberName, CastMemberRole, CastMemberGender
FROM CastMember;

-- name: UpdateCastMember :exec
UPDATE CastMember
SET CastMemberName = $2, CastMemberRole = $3, CastMemberGender = $4
WHERE CastMemberId = $1;

-- name: DeleteCastMember :exec
DELETE FROM CastMember
WHERE CastMemberId = $1;

-- name: CreateFilmCast :exec
INSERT INTO FilmCast(FIlmId, CastMemberId)
VALUES($1, $2)
RETURNING FIlmId, CastMemberId;

-- name: GetFilmCast :one
SELECT FIlmId, CastMemberId
FROM FilmCast
WHERE FIlmId = $1 AND CastMemberId = $2;

-- name: GetFilmCasts :many
SELECT FIlmId, CastMemberId
FROM FilmCast;

-- name: DeleteFilmCast :exec
DELETE FROM FilmCast
WHERE FIlmId = $1 AND CastMemberId = $2;

-- name: CreateFilmGenre :exec
INSERT INTO FilmGenre(Genre_GenreID, Film_FilmId)
VALUES($1, $2)
RETURNING Genre_GenreID, Film_FilmId;

-- name: GetFilmGenre :one
SELECT Genre_GenreID, Film_FilmId
FROM FilmGenre
WHERE Genre_GenreID = $1 AND Film_FilmId = $2;

-- name: GetFilmGenres :many
SELECT Genre_GenreID, Film_FilmId
FROM FilmGenre;

-- name: DeleteFilmGenre :exec
DELETE FROM FilmGenre
WHERE Genre_GenreID = $1 AND Film_FilmId = $2;

-- name: CreateFilmUserStatus :exec
INSERT INTO FilmUserStatus(Status_StatusId, Film_FilmId, User_UserID)
VALUES($1, $2, $3)
RETURNING Status_StatusId, Film_FilmId, User_UserID;

-- name: GetFilmUserStatus :one
SELECT Status_StatusId, Film_FilmId, User_UserID
FROM FilmUserStatus
WHERE Film_FilmId = $1 AND User_UserID = $2;

-- name: GetFilmUserStatuses :many
SELECT Status_StatusId, Film_FilmId, User_UserID
FROM FilmUserStatus;

-- name: UpdateFilmUserStatus :exec
UPDATE FilmUserStatus
SET Status_StatusId = $3
WHERE Film_FilmId = $1 AND User_UserID = $2;

-- name: DeleteFilmUserStatus :exec
DELETE FROM FilmUserStatus
WHERE Film_FilmId = $1 AND User_UserID = $2;

-- name: CreateGenre :exec
INSERT INTO Genre(GenreID, GenreName)
VALUES($1, $2)
RETURNING GenreID, GenreName;

-- name: GetGenre :one
SELECT GenreID, GenreName
FROM Genre
WHERE GenreID = $1;

-- name: GetGenres :many
SELECT GenreID, GenreName
FROM Genre;

-- name: UpdateGenre :exec
UPDATE Genre
SET GenreName = $2
WHERE GenreID = $1;

-- name: DeleteGenre :exec
DELETE FROM Genre
WHERE GenreID = $1;

-- name: CreateReview :exec
INSERT INTO Review(ReviewStars, ReviewDescription, ReviewUserID, ReviewFilmId)
VALUES($1, $2, $3, $4)
RETURNING ReviewStars, ReviewDescription, ReviewUserID, ReviewFilmId;

-- name: GetReview :one
SELECT ReviewStars, ReviewDescription, ReviewUserID, ReviewFilmId
FROM Review
WHERE ReviewFilmId = $1 AND ReviewUserID = $2;

-- name: GetReviews :many
SELECT ReviewStars, ReviewDescription, ReviewUserID, ReviewFilmId
FROM Review;

-- name: UpdateReview :exec
UPDATE Review
SET ReviewStars = $3, ReviewDescription = $4
WHERE ReviewFilmId = $1 AND ReviewUserID = $2;

-- name: DeleteReview :exec
DELETE FROM Review
WHERE ReviewFilmId = $1 AND ReviewUserID = $2;

-- name: CreateStatus :exec
INSERT INTO Status(StatusId, StatusName)
VALUES($1, $2)
RETURNING StatusId, StatusName;

-- name: GetStatus :one
SELECT StatusId, StatusName
FROM Status
WHERE StatusId = $1;

-- name: GetStatuses :many
SELECT StatusId, StatusName
FROM Status;

-- name: UpdateStatus :exec
UPDATE Status
SET StatusName = $2
WHERE StatusId = $1;

-- name: DeleteStatus :exec
DELETE FROM Status
WHERE StatusId = $1;

-- name: CreateUser :exec
INSERT INTO "User"(UserID, UserName, UserPasswordHashed, UserBirthday, UserGenre)
VALUES($1, $2, $3, $4, $5)
RETURNING UserID, UserName, UserPasswordHashed, UserBirthday, UserGenre;

-- name: GetUser :one
SELECT UserID, UserName, UserPasswordHashed, UserBirthday, UserGenre
FROM "User"
WHERE UserID = $1;

-- name: GetUsers :many
SELECT UserID, UserName, UserPasswordHashed, UserBirthday, UserGenre
FROM "User";

-- name: UpdateUser :exec
UPDATE "User"
SET UserName = $2, UserPasswordHashed = $3, UserBirthday = $4, UserGenre = $5
WHERE UserID = $1;

-- name: DeleteUser :exec
DELETE FROM "User"
WHERE UserID = $1;