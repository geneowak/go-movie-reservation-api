-- name: CreateMovie :one
INSERT INTO
    movies(
        id,
        name,
        description,
        duration_in_mins,
        trailer_url,
        genre,
        pg_rating,
        experience_types,
        created_at,
        updated_at
    )
VALUES
    (
        uuidv7(),
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        NOW(),
        NOW()
    )
RETURNING
    *;

-- name: GetMovieById :one
SELECT
    *
FROM
    movies
WHERE
    id = $1;

-- name: CheckMovieById :one
SELECT
    EXISTS(
        SELECT
            1
        FROM
            movies
        WHERE
            id = $1
    );

-- name: GetShowingMovies :many
SELECT
    movies.*,
    coalesce(
        (
            SELECT
                jsonb_agg(
                    to_jsonb(st)
                    ORDER BY
                        st.start_date,
                        st.start_time
                )
            FROM
                show_times st
            WHERE
                st.movie_id = movies.id
                AND st.end_date >= NOW()
        ),
        '[]'::jsonb
    ) AS show_times
FROM
    movies
WHERE
    EXISTS (
        SELECT
            1
        FROM
            show_times st
        WHERE
            st.movie_id = movies.id
            AND st.end_date >= NOW()
    )
ORDER BY
    (
        SELECT
            MIN(st.start_date)
        FROM
            show_times st
        WHERE
            st.movie_id = movies.id
            AND st.end_date >= NOW()
    ) ASC;
