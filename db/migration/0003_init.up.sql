insert into prediapp.prode_sessions (user_id, session_id, p1, p2, p3, created_at, updated_at)
VALUES 
(3, 1, 1, 2, 3, NOW(), NOW()),
(1, 1, 1, 2, 3, NOW(), NOW());

insert into prediapp.prode_carreras (user_id , session_id, p1, p2, p3, p4, p5, vsc, sc, dnf, created_at, updated_at)
values 
(3, 2, 1, 2, 3, 4, 5, true, true, 2, NOW(), NOW()),
(1, 2, 1, 2, 3, 4, 5, true, true, 2, NOW(), NOW());