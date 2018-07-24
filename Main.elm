-- Read more about this program in the official Elm guide:
-- https://guide.elm-lang.org/architecture/effects/http.html


module Main exposing (..)

import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import Json.Decode
import Time exposing (Time, second)
import Round
import Dict exposing (Dict)


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
    ( { status = { nodes = [], measurements = Dict.fromList [] }, error = "" }, fetchResults )


subscriptions : Model -> Sub Msg
subscriptions model =
    Time.every (1 * second) Tick


type alias Model =
    { status : Status
    , error : String
    }


type alias Status =
    { nodes : List Node
    , measurements : Dict String (Dict String Measurement)
    }


type alias Node =
    { hostName : String
    , hostIP : String
    , podName : String
    , podIP : String
    }


type alias Measurement =
    { delay : Int
    , timestamp : Int
    , error : String
    }


type Msg
    = Tick Time
    | NewResults (Result Http.Error Status)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        Tick newTime ->
            ( model, fetchResults )

        NewResults (Ok newResults) ->
            ( { model | status = newResults, error = "" }, Cmd.none )

        NewResults (Err error) ->
            ( { model | error = toString error }, Cmd.none )


view : Model -> Html Msg
view model =
    div [ class "goldpinger" ]
        [ css "https://necolas.github.io/normalize.css/8.0.0/normalize.css"
        , css "styles.css"
        , Html.h1 [] [ text "Goldpinger" ]
        , viewTable model.status
        , printError model.error
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


viewTable : Status -> Html Msg
viewTable status =
    let
        headline =
            viewHosts status.nodes

        rows =
            viewRows status
    in
        Html.table [] (List.concat [ [ headline ], rows ])


viewHosts : List Node -> Html Msg
viewHosts nodes =
    Html.tr []
        (List.concat
            [ [ Html.td [] [] ]
            , List.map (\node -> Html.td [] [ div [ class "to" ] [ text "to ", text node.hostName ] ]) nodes
            ]
        )


viewRows : Status -> List (Html Msg)
viewRows status =
    List.map (\h -> viewRow h status.nodes status.measurements) status.nodes


viewRow : Node -> List Node -> Dict String (Dict String Measurement) -> Html Msg
viewRow node nodes measurements =
    let
        m =
            Dict.get node.hostName measurements
    in
        case m of
            Nothing ->
                Html.tr []
                    (List.concat
                        [ [ Html.td [] [ text "from ", text node.hostName ] ]
                        , List.map (\node -> emptyCell) nodes
                        ]
                    )

            Just m ->
                Html.tr []
                    (List.concat
                        [ [ Html.td [] [ text "from ", text node.hostName ] ]
                        , List.map (\node -> viewCell node.hostName m) nodes
                        ]
                    )


viewCell : String -> Dict String Measurement -> Html Msg
viewCell target measurements =
    let
        measurement =
            Dict.get target measurements
    in
        case measurement of
            Nothing ->
                emptyCell

            Just m ->
                printColored m


printColored : Measurement -> Html Msg
printColored ping =
    let
        delayInMilisec =
            toFloat ping.delay / 100000

        display =
            Round.round 2 delayInMilisec
    in
        if delayInMilisec > 50 then
            Html.td [ class "high ping" ] [ text display ]
        else if delayInMilisec > 25 then
            Html.td [ class "med ping" ] [ text display ]
        else
            Html.td [ class "low ping" ] [ text display ]


emptyCell : Html Msg
emptyCell =
    Html.td [ class "empty ping" ] []


fetchResults : Cmd Msg
fetchResults =
    Http.send NewResults (Http.get "./status.json" decodeStatus)


decodeStatus : Json.Decode.Decoder Status
decodeStatus =
    Json.Decode.map2 Status
        (Json.Decode.field "nodes" (Json.Decode.list decodeNode))
        (Json.Decode.field "measurements" (Json.Decode.dict (Json.Decode.dict decodeMeasurement)))


decodeNode : Json.Decode.Decoder Node
decodeNode =
    Json.Decode.map4 Node
        (Json.Decode.field "hostName" Json.Decode.string)
        (Json.Decode.field "hostIP" Json.Decode.string)
        (Json.Decode.field "podName" Json.Decode.string)
        (Json.Decode.field "podIP" Json.Decode.string)


decodeMeasurement : Json.Decode.Decoder Measurement
decodeMeasurement =
    Json.Decode.map3 Measurement
        (Json.Decode.field "delay" Json.Decode.int)
        (Json.Decode.field "timestamp" Json.Decode.int)
        (Json.Decode.field "error" Json.Decode.string)
