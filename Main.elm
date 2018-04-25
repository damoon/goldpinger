-- Read more about this program in the official Elm guide:
-- https://guide.elm-lang.org/architecture/effects/http.html


module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Http
import Json.Decode
import Time exposing (Time, second)
import Debug exposing (log)


main : Program Never Model Msg
main =
    Html.program
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }


init : ( Model, Cmd Msg )
init =
    ( { hosts = [], error = "" }
    , fetchResults
    )


subscriptions : Model -> Sub Msg
subscriptions model =
    Time.every (1 * second) Tick


type alias Model =
    { hosts : List Host
    , error : String
    }


type alias Host =
    { source : String
    , pings : List Ping
    }


type alias Ping =
    { target : String
    , delay : Int
    , timestamp : Int
    }


type Msg
    = Tick Time
    | NewResults (Result Http.Error (List Host))


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Tick newTime ->
            ( model, fetchResults )

        NewResults (Ok newResults) ->
            ( { model | hosts = newResults, error = "" }, Cmd.none )

        NewResults (Err error) ->
            ( { model | error = toString error }, Cmd.none )


view : Model -> Html Msg
view model =
    div []
        [ printError model.error, viewTable model.hosts ]


printError : String -> Html Msg
printError error =
    if error == "" then
        text ""
    else
        div [ Html.Attributes.class "error" ] [ text error ]


viewTable : List Host -> Html Msg
viewTable hosts =
    Html.table []
        (List.map
            viewRow
            hosts
        )


viewRow : Host -> Html Msg
viewRow host =
    Html.tr []
        (List.concat
            [ [ Html.td [] [ Html.text host.source ] ]
            , (List.map
                viewPing
                host.pings
              )
            ]
        )


viewPing : Ping -> Html Msg
viewPing ping =
    Html.td []
        [ text "target"
        , text ping.target
        , br [] []
        , text "delay"
        , text (toString ping.delay)
        , br [] []
        , text "timestamp"
        , text (toString ping.timestamp)
        ]


fetchResults : Cmd Msg
fetchResults =
    let
        url =
            "http://localhost:8080/"
    in
        Http.send NewResults (Http.get url (Json.Decode.list decodeHost))


decodeHost : Json.Decode.Decoder Host
decodeHost =
    Json.Decode.map2 Host
        (Json.Decode.field "source" Json.Decode.string)
        (Json.Decode.field "pings" (Json.Decode.list decodePing))


decodePing : Json.Decode.Decoder Ping
decodePing =
    Json.Decode.map3 Ping
        (Json.Decode.field "target" Json.Decode.string)
        (Json.Decode.field "delay" Json.Decode.int)
        (Json.Decode.field "timestamp" Json.Decode.int)
