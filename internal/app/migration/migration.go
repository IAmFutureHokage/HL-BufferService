package migration

const CreateTablesTelegramAndPhenomenia = `
CREATE TABLE IF NOT EXISTS telegram (
    id TEXT PRIMARY KEY,
    groupId TEXT,
    telegramCode TEXT,
    postCode TEXT,
    dateTime TIMESTAMP,
    endBlockNum SMALLINT,
    isDangerous BOOLEAN,
    waterLevelOnTime INTEGER,
    deltaWaterLevel INTEGER,
    waterLevelOn20h INTEGER,
    waterTemperature DOUBLE PRECISION,
    airTemperature INTEGER,
    icePhenomeniaState SMALLINT,
    ice SMALLINT,
    snow SMALLINT,
    waterflow DOUBLE PRECISION,
    precipitationValue DOUBLE PRECISION,
    precipitationDuration SMALLINT,
    reservoirDate TIMESTAMP,
    headwaterLevel INTEGER,
    averageReservoirLevel INTEGER,
    downstreamLevel INTEGER,
    reservoirVolume DOUBLE PRECISION,
    isReservoirWaterInflowDate TIMESTAMP,
    inflow DOUBLE PRECISION,
    reset DOUBLE PRECISION
);

CREATE TABLE IF NOT EXISTS phenomenia (
    iId TEXT PRIMARY KEY,
    telegramId TEXT REFERENCES telegram(id),
    phenomen SMALLINT,
    isUntensity BOOLEAN,
    intensity SMALLINT
);
`
