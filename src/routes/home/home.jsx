import { useEffect, useState } from "react";
import { getTime, getWeatherInfo } from "./home-fns";
import "./home.scss";

export const Home = ({ setError }) => {
    const [date, setDate] = useState("");
    const [time, setTime] = useState("");
    const [weatherInfo, setWeatherInfo] = useState({});
    const [tempCelsius, setTempCelsius] = useState("");
    const [tempFahrenheit, setTempFahrenheit] = useState("");
    const [feelsLikeCelsius, setFeelsLikeCelsius] = useState("");
    const [feelsLikeFahrenheit, setFeelsLikeFahrenheit] = useState("");
    const [isCelsius, setIsCelsius] = useState(true);
    const [timeZone, setTimeZone] = useState("");

    setError(false);


    useEffect(() => {
        document.title = `Lesi | Home`;

    }, []);

    useEffect(() => {
        let dateTimeIntervalId;
        let weatherIntervalId;

        async function updateDateTime() {
            const data = await getTime();
            setTime(data[0]);
            setDate(data[1]);
        }

        async function updateWeatherInfo() {
            const data = await getWeatherInfo();
            const json = JSON.parse(data);
            // console.log(json);
            setWeatherInfo(json);
            setTempCelsius(json?.current?.temp_c);
            setTempFahrenheit(json?.current?.temp_f);
            setFeelsLikeCelsius(json?.current?.feelslike_c);
            setFeelsLikeFahrenheit(json?.current?.feelslike_f);
        }

        updateDateTime();
        updateWeatherInfo();

        dateTimeIntervalId = setInterval(updateDateTime, 1000);
        weatherIntervalId = setInterval(updateWeatherInfo, 300000);

        return () => {
            clearInterval(dateTimeIntervalId);
            clearInterval(weatherIntervalId);
        };
    }, []);

    // function convertToFahrenheit(temp) {
    //     return (parseFloat(temp) * 9 / 5) + 32;
    // }

    const toggleTemperatureUnit = () => {
        setIsCelsius(!isCelsius);
    };

    return (
        <main>
            <div className="home">
                <div className="description">
                    <h1>{time}</h1>
                    <h2>{date}</h2>
                </div>

                <div className="weather">
                    <div className="values">
                        Current weather in {weatherInfo?.location?.name}
                    </div>
                    <div className="values clickable"
                        onClick={toggleTemperatureUnit}>
                        Temperature: {isCelsius ? `${tempCelsius}°C` : `${tempFahrenheit}°F`}
                    </div>
                    <div className="values clickable"
                        onClick={toggleTemperatureUnit}>
                        Feels Like: {isCelsius ? `${feelsLikeCelsius}°C` : `${feelsLikeFahrenheit}°F`}
                    </div>
                    <div className="values">
                        Humidity: {weatherInfo?.current?.humidity}&#37;
                    </div>
                    <div className="values">
                        Wind Speed: {weatherInfo?.current?.wind_mph}&#37;
                    </div>
                </div>
            </div>





        </main>
    )
}