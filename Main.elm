-- Read more about this program in the official Elm guide:
-- https://guide.elm-lang.org/architecture/effects/http.html


module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import Json.Decode
import Time exposing (Time, second)
import Round


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
    { hosts : List Source
    , error : String
    }


type alias Source =
    { hostName : String
    , measurements : List Measurement
    }


type alias Measurement =
    { target : String
    , delay : Int
    , timestamp : Int
    , error : String
    }


type Msg
    = Tick Time
    | NewResults (Result Http.Error (List Source))


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
    div [ class "goldpinger" ]
        [ css "https://necolas.github.io/normalize.css/8.0.0/normalize.css"
        , css "styles.css"
        , printError model.error
        , Html.h1 [] [ text "Goldpinger" ]
        , viewTable model.hosts
        ]


css : String -> Html Msg
css path =
    node "link" [ rel "stylesheet", href path ] []


printError : String -> Html Msg
printError error =
    if error == "" then
        text ""
    else
        div [ class "error" ] [ text error ]


viewTable : List Source -> Html Msg
viewTable hosts =
    let
        sorted =
            List.sortBy .hostName hosts
    in
        Html.table [] (List.concat [ [ viewHosts sorted ], List.map viewRow sorted ])


viewHosts : List Source -> Html Msg
viewHosts hosts =
    Html.tr []
        (List.concat
            [ [ Html.td [] [] ]
            , List.map (\h -> Html.td [] [ div [ class "to" ] [ text "to ", text h.hostName ] ]) hosts
            ]
        )


viewRow : Source -> Html Msg
viewRow host =
    let
        sorted =
            List.sortBy .target host.measurements
    in
        Html.tr []
            (List.concat
                [ [ Html.td [] [ text "from ", text host.hostName ] ]
                , List.map (viewPing host.hostName) sorted
                ]
            )


viewPing : String -> Measurement -> Html Msg
viewPing source ping =
    if ping.target == source then
        Html.td [ Html.Attributes.class "local-machine" ] [ text "-" ]
    else
        printColored ping.delay


printColored : Int -> Html Msg
printColored delay =
    let
        delayInMilisec =
            toFloat delay / 100000

        display =
            Round.round 2 delayInMilisec
    in
        if delayInMilisec > 50 then
            Html.td [ class "high ping" ] [ text display ]
        else if delayInMilisec > 25 then
            Html.td [ class "med ping" ] [ text display ]
        else
            Html.td [ class "low ping" ] [ text display ]


fetchResults : Cmd Msg
fetchResults =
    Http.send NewResults (Http.get "./status.json" (Json.Decode.list decodeHost))


decodeHost : Json.Decode.Decoder Source
decodeHost =
    Json.Decode.map2 Source
        (Json.Decode.field "hostName" Json.Decode.string)
        (Json.Decode.field "measurements" (Json.Decode.list decodePing))


decodePing : Json.Decode.Decoder Measurement
decodePing =
    Json.Decode.map4 Measurement
        (Json.Decode.field "target" Json.Decode.string)
        (Json.Decode.field "delay" Json.Decode.int)
        (Json.Decode.field "timestamp" Json.Decode.int)
        (Json.Decode.field "error" Json.Decode.string)
