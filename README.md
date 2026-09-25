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
- [x] GET /api/show-times/{showTimeId}
- [x] GET /api/locations
- [x] GET /api/locations/{locationId}
- [x] PUT /api/cinemas/{cinemaId} #e.g updating seat map
- [x] GET /api/admin/movies (returns all the movies for admin management)
- [ ] PUT /api/movies/{movieId}
- [ ] DELETE /api/movies/{movieId}
- [ ] PUT /api/users/{userId} # assigning or removing admins

USER ROUTES
- [x] GET /api/movies?genre={genre}&show-time={time}&date={date} (gets currently showing movies)
- [x] GET /api/movies/{movieId} ->movie with show-times
- [x] POST /api/seats/reserve -> reserve seats for 5mins while browsing
- [x] POST /api/seats/book -> confirm seat booking
- [x] GET /api/bookings
- [x] DELETE /api/bookings/{reservationId} -> user can remove seats they have booked


TODO (clean up after initial implementation)
- [ ] validate the reserveSeatRequest, the user should own all the reservations and they should not be already booked


## Postmoterm
If I was restarting this project or implementing it for a business, here are some of the things I would change given what I have learnt from it
- I would rename the `reservations` table to `bookings` table and have the possible status as `pending` and `booked`. I would keep the same structure. I think naming the table bookings is more intuitive because then the routes would all have the prefix instead of mixing reservations and bookings  `/api/bookings`
- I would change all the `timestamp` fields to `timestampstz` with the timezone of UTC. This would avoid time related bugs like the one I had when trying to calculate how much time had passed since a seat was reserved. I put it a fix for this in the commit `7a47fd65cc464d79491224454deac8bc58422da7` but that would not have been needed if I had used timezones from the start.
