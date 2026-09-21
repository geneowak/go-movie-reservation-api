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
        AND (
            sqlc.narg(time)::time IS NULL
            OR st.start_time >= sqlc.narg(time)
        )
        AND (
            sqlc.narg(date)::timestamp IS NULL
            OR st.start_date <= sqlc.narg(date)
        )
    GROUP BY
        st.movie_id
)
SELECT
    movies.*,
    fs.show_times
FROM
    movies
    INNER JOIN filtered_shows fs ON movies.id = fs.movie_id
WHERE
    (
        sqlc.narg(genre)::text IS NULL
        OR movies.genre::jsonb @> jsonb_build_array(sqlc.narg(genre)::text)
    )
ORDER BY
    fs.first_show_date ASC;
