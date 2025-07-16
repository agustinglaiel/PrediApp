CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255),
    score INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_email_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP NULL,
    phone_number VARCHAR(20),
    provider VARCHAR(255),
    provider_id VARCHAR(255),
    imagen_perfil MEDIUMBLOB,
    imagen_mime_type VARCHAR(50),
    UNIQUE INDEX idx_username (username),
    UNIQUE INDEX idx_email (email),
    INDEX idx_deleted_at (deleted_at)
);

CREATE TABLE drivers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    broadcast_name VARCHAR(100),
    country_code VARCHAR(10),
    driver_number INT,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    full_name VARCHAR(100),
    name_acronym VARCHAR(10),
    headshot_url VARCHAR(200),
    team_name VARCHAR(100),
    activo BOOLEAN,
    INDEX idx_driver_name (first_name, last_name)
);

CREATE TABLE sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    weekend_id INT,
    circuit_key INT,
    circuit_short_name VARCHAR(255),
    country_code VARCHAR(255),
    country_key INT,
    country_name VARCHAR(255),
    location VARCHAR(255),
    session_key INT,
    session_name VARCHAR(255),
    session_type VARCHAR(255),
    date_start TIMESTAMP,
    date_end TIMESTAMP,
    year INT,
    vsc BOOLEAN,
    sf BOOLEAN,
    dnf INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_deleted_at (deleted_at)
);

CREATE TABLE results (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id INT,
    driver_id INT,
    position INT,
    fastest_lap_time DOUBLE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status LONGTEXT,
    INDEX idx_session_id (session_id),
    INDEX idx_driver_id (driver_id),
    CONSTRAINT fk_results_session FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_results_driver FOREIGN KEY (driver_id) REFERENCES drivers(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE prode_carreras (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    session_id INT NOT NULL,
    p1 INT,
    p2 INT,
    p3 INT,
    p4 INT,
    p5 INT,
    vsc BOOLEAN,
    sc BOOLEAN,
    dnf INT,
    score INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_prode_carreras_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_prode_carreras_session FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_prode_carreras_p1 FOREIGN KEY (p1) REFERENCES drivers(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_prode_carreras_p2 FOREIGN KEY (p2) REFERENCES drivers(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_prode_carreras_p3 FOREIGN KEY (p3) REFERENCES drivers(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_prode_carreras_p4 FOREIGN KEY (p4) REFERENCES drivers(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_prode_carreras_p5 FOREIGN KEY (p5) REFERENCES drivers(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE prode_sessions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    session_id INT NOT NULL,
    p1 INT,
    p2 INT,
    p3 INT,
    score INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_prode_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_prode_sessions_session FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_prode_sessions_p1 FOREIGN KEY (p1) REFERENCES drivers(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_prode_sessions_p2 FOREIGN KEY (p2) REFERENCES drivers(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_prode_sessions_p3 FOREIGN KEY (p3) REFERENCES drivers(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE TABLE `groups` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    group_code VARCHAR(8) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE INDEX idx_group_code (group_code),
    INDEX idx_deleted_at (deleted_at)
);

CREATE TABLE group_x_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_id INT NOT NULL,
    user_id INT NOT NULL,
    group_role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_group_id (group_id),
    INDEX idx_user_id (user_id),
    CONSTRAINT fk_group_x_users_group FOREIGN KEY (group_id) REFERENCES `groups`(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_group_x_users_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    parent_post_id INT,
    body VARCHAR(500) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_parent_post_id (parent_post_id),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at),
    CONSTRAINT fk_posts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT fk_posts_parent_post FOREIGN KEY (parent_post_id) REFERENCES posts(id)
);