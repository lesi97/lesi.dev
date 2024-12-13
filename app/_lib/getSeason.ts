function getSeason() {
    const date = new Date();
    const month = date.getMonth() + 1;
    const day = date.getDate();

    switch (true) {
        case month === 2 && day >= 7 && day <= 14:
            return 'Valentines';
        case month === 10 && day >= 15 && day <= 31:
            return 'Halloween';
        case month === 11 && day === 1:
            return 'Lesi-Birthday';
        case month === 12 && day >= 1 && day <= 26:
            return 'Christmas';
        case month === 12 && day >= 27 && day <= 31:
        case month === 1 && day === 1:
            return 'New-Years';
        default:
            return null;
    }
}

export default getSeason;
