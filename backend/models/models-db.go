package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (r *DBModel) GetMovieDB(id int) (*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, title, description, year, release_date, runtime, mpaa_rating, created_at, updated_at, coalesce(poster, '') from movies where id = $1`
	row := r.DB.QueryRowContext(ctx, query, id)
	var movie Movie
	err := row.Scan(&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.Year,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.Poster,
	)
	if err != nil {
		return nil, err
	}

	query = `select mg.id, mg.movie_id, mg.genre_id, g.genre_name from movies_genres mg left join genres g on (g.id = mg.genre_id) where mg.movie_id = $1`
	rows, _ := r.DB.QueryContext(ctx, query, id)
	defer rows.Close()

	genres := make(map[int]string)
	for rows.Next() {
		var mg MovieGenre
		err := rows.Scan(
			&mg.ID,
			&mg.MovieID,
			&mg.GenreID,
			&mg.Genre.Name,
		)
		if err != nil {
			return nil, err
		}
		genres[mg.ID] = mg.Genre.Name
	}
	movie.MovieGenre = genres

	return &movie, nil
}

func (r *DBModel) InsertMovieDB(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `insert into movies (title, description, year, release_date, runtime, mpaa_rating, created_at, updated_at, poster) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.DB.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.CreatedAt,
		movie.UpdatedAt,
		movie.Poster,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DBModel) EditMovieDB(movie Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update movies set title=$1, description=$2, year=$3, release_date=$4, runtime=$5, mpaa_rating=$6, updated_at=$7, poster=$8 where id=$9`

	_, err := r.DB.ExecContext(ctx, query,
		movie.Title,
		movie.Description,
		movie.Year,
		movie.ReleaseDate,
		movie.Runtime,
		movie.Rating,
		movie.UpdatedAt,
		movie.Poster,
		movie.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DBModel) DeleteMovieDB(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `delete from movies where id = $1`

	_, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *DBModel) GetMoviesDB(genre ...int) ([]*Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(genre) > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", genre[0])
	}

	query := fmt.Sprintf(`select id, title, description, year, release_date, runtime, mpaa_rating, created_at, updated_at, coalesce(poster, '') from movies %s order by title`, where)
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.Year,
			&movie.ReleaseDate,
			&movie.Runtime,
			&movie.Rating,
			&movie.CreatedAt,
			&movie.UpdatedAt,
			&movie.Poster,
		)
		if err != nil {
			return nil, err
		}

		genreQuery := `select mg.id, mg.movie_id, mg.genre_id, g.genre_name from movies_genres mg left join genres g on (g.id = mg.genre_id) where mg.movie_id = $1`
		rows, _ := r.DB.QueryContext(ctx, genreQuery, movie.ID)

		genres := make(map[int]string)
		for rows.Next() {
			var mg MovieGenre
			err := rows.Scan(
				&mg.ID,
				&mg.MovieID,
				&mg.GenreID,
				&mg.Genre.Name,
			)
			if err != nil {
				return nil, err
			}
			genres[mg.ID] = mg.Genre.Name
		}
		movie.MovieGenre = genres
		movies = append(movies, &movie)
	}

	return movies, nil
}

func (r *DBModel) GetGenresDB() ([]*Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, genre_name, created_at, updated_at from genres order by genre_name`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*Genre
	for rows.Next() {
		var mg Genre
		err := rows.Scan(
			&mg.ID,
			&mg.Name,
			&mg.CreatedAt,
			&mg.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &mg)
	}

	return genres, nil
}

func (r *DBModel) GetUser(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, email, pass from users where email = $1`
	row := r.DB.QueryRowContext(ctx, query, email)
	var user User
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
