
async function updateCounter() {
    let counter = await fetch("./php/homePageCounter");
    let counterResponse = await counter.text();
    console.log(counterResponse);
}
updateCounter();

export async function getTime() {

    let timeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;

    let time = await fetch(`./php/time?timezone=${timeZone}`, {
        method: 'GET'
    });
    let timeResponse = await time.text();

    const daysOfWeek = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
    const monthsOfYear = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];

    const currentDate = new Date();

    const dayName = daysOfWeek[currentDate.getDay()];
    const dayNumber = currentDate.getDate();
    const monthName = monthsOfYear[currentDate.getMonth()];
    const year = currentDate.getFullYear();
    const ordinal = getSuffix(dayNumber);

    let dateResponse = `${dayName} ${dayNumber}${ordinal} ${monthName} ${year}`;

    return [timeResponse, dateResponse];

}

function getSuffix(number) {
    if (number % 10 == 1 && number != 11) {
        return 'st';
    } else if (number % 10 == 2 && number != 12) {
        return 'nd';
    } else if (number % 10 == 3 && number != 13) {
        return 'rd';
    } else {
        return 'th';
    }
}

async function getGeoLocationData() {
    try {
        let response = await fetch("./php/geolocationApi", {
            method: "GET"
        });
        let data = await response.json();
        if (data.status === "success") {
            return {
                latitude: data.lat,
                longitude: data.lon
            };
        }
    } catch (error) {
        console.error('There was an error fetching the weather information:', error);
    }
}

export async function getWeatherInfo() {
    try {
        const { latitude, longitude } = await getGeoLocationData();
        let response = await fetch("./php/weatherApi?latitude=" + latitude + "&longitude=" + longitude, {
            method: "GET"
        });
        let data = await response.text();
        return data;
    } catch (error) {
        console.error('There was an error fetching the weather information:', error);
    }
}