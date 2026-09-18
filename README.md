## Movie Reservation System API

## STATUS: INPROGRESS

Project details can be found [here](https://roadmap.sh/projects/movie-reservation-system)

## Draft

This readme is going to be updated with more details on the implementation approach taken but for now the above reference will do

### APIs to be implemented

- [x] POST /api/signUp
- [x] POST /api/login
- [x] POST /api/refresh (auth tokens)
- [x] POST /api/revoke (refresh tokens)

ADMIN ROUTES
- [x] POST /api/locations
- [x] POST /api/cinemas
- [x] POST /api/movies
- [x] POST /api/movies/{movieId}/show-times
- [ ] GET /api/show-times
- [ ] GET /api/locations
- [ ] GET /api/cinemas
- [ ] PUT /api/cinemas/{cinemaId} #e.g updating seat map
- [ ] GET /api/admin/movies (returns all the movies for admin management)
- [ ] PUT /api/movies/{movieId}
- [ ] DELETE /api/movies/{movieId}
- [ ] PUT /api/users/{userId} # assigning or removing admins

USER ROUTES
- [ ] GET /api/movies?genre={genre}&show-time={time}&date={date} (gets currently showing movies)
- [ ] GET /api/movies/{movieId} ->movie with show-times
- [ ] POST /api/seats/reserve -> reserve seats for 5mins while browsing
- [ ] POST /api/seats/book -> confirm seat booking
- [ ] GET /api/bookings
- [ ] PUT /api/bookings/{bookingId} -> user can add/remove seats they have booked
