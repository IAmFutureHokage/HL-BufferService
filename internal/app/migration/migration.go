package migration

const CreateTablesTelegramAndPhenomenia = `
CREATE TABLE IF NOT EXISTS telegram (
    id TEXT PRIMARY KEY,
    groupid TEXT,
    telegramcode TEXT,
    postcode TEXT,
    datetime TIMESTAMPTZ,
    endblocknum SMALLINT,
    isdangerous BOOLEAN,
    waterlevelontime INTEGER,
    deltaWaterlevel INTEGER,
    waterlevelon20h INTEGER,
    watertemperature DOUBLE PRECISION,
    airtemperature INTEGER,
    icephenomeniastate SMALLINT,
    ice INTEGER,
    snow SMALLINT,
    waterflow DOUBLE PRECISION,
    precipitationvalue DOUBLE PRECISION,
    precipitationduration SMALLINT,
    reservoirdate TIMESTAMPTZ,
    headwaterlevel INTEGER,
    averagereservoirlevel INTEGER,
    downstreamlevel INTEGER,
    reservoirvolume DOUBLE PRECISION,
    isreservoirwaterinflowdate TIMESTAMPTZ,
    inflow DOUBLE PRECISION,
    reset DOUBLE PRECISION
);

CREATE TABLE IF NOT EXISTS phenomenia (
    id TEXT PRIMARY KEY,
    telegramid TEXT REFERENCES telegram(id),
    phenomen SMALLINT,
    isuntensity BOOLEAN,
    intensity SMALLINT
);
`
