import React, {Fragment, useEffect, useState} from "react";
import {Link, useLocation, useParams} from "react-router-dom";

export default function Genre() {
    const params = useParams()
    const location = useLocation()

    let [movies, setMovies] = useState([])
    let [genreName, setGenreName] = useState("")
    const [error, setError] = useState(null)

    useEffect(() => {
        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/movies/` + params.id)
            .then((response) => {
                if (response.status !== 200)
                    setError("Invalid response code: ", response.status)
                else
                    setError(null)
                return response.json()
            })
            .then((json) => {
                setMovies(json.movies)
                setGenreName(location.state.genreName)
            })
    }, [params.id, location.state.genreName])

    if (!movies)
        movies = []

    if (error !== null)
        return <div>Error: {error.message}</div>
    else {
        return (
            <Fragment>
                <h2>Genre: {genreName}</h2>
                <div className={"list-group"}>
                    {movies.map((m, index) => (
                        <Link to={`/movies/${m.id}`}
                              className={"list-group-item list-group-item-action"} key={index}>{m.title}</Link>
                    ))}
                </div>
            </Fragment>
        )
    }
}