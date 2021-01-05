CREATE TABLE flights(
                             flight_id BIGSERIAL PRIMARY KEY,
                             fromCity TEXT NOT NULL,
                             toCity TEXT NOT NULL,
                             duration BIGINT NOT NULL,
                             cost BIGINT NOT NULL,
                             departureTime DATE NOT NULL DEFAULT CURRENT_DATE
);
