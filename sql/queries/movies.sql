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

-- name: GetAllMovies :many
SELECT
    movies.*,
    COALESCE(
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
        ),
        '[]'::jsonb
    ) AS show_times
FROM
    movies
ORDER BY
    (
        SELECT
            min(st.start_date)
        FROM
            show_times st
        WHERE
            st.movie_id = movies.id
    ) ASC;

-- name: GetShowingMovies :many
WITH filtered_shows AS (
    SELECT
        st.movie_id,
        jsonb_agg(
            to_jsonb(st)
            ORDER BY
                st.start_date,
                st.start_time
        ) AS show_times,
        min(st.start_date) AS first_show_date
    FROM
        show_times st
    WHERE
        st.end_date >= CURRENT_DATE
    GROUP BY
        st.movie_id
)
SELECT
    movies.*,
    fs.show_times
FROM
    movies
    INNER JOIN filtered_shows fs ON movies.id = fs.movie_id
ORDER BY
    fs.first_show_date ASC;
