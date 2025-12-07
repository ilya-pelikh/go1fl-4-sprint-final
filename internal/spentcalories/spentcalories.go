package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splittedStrings := strings.Split(data, ",")

	if len(splittedStrings) != 3 {
		return 0, "", 0, errors.New("недостаточно данных")
	}

	steps, err := strconv.Atoi(splittedStrings[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("некорректные данные шагов")
	}

	activity := splittedStrings[1]

	duration, err := time.ParseDuration(splittedStrings[2])
	if err != nil || duration <= 0 {
		return 0, "", 0, errors.New("длительность меньше либо равна нулю")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	lenOfStep := stepLengthCoefficient * height

	distanceMeters := lenOfStep * float64(steps)

	return distanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}

	dist := distance(steps, height)

	return dist / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0.0, errors.New("данные не подходят")
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0.0, errors.New("данные не подходят")
	}

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0.0, errors.New("данные не подходят")
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0.0, errors.New("данные не подходят")
	}

	durationMinutes := duration.Minutes()

	baseCalories := (weight * speed * durationMinutes) / minInH

	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)

	speed := meanSpeed(steps, height, duration)

	var calories float64
	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	hours := duration.Hours()

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		hours,
		dist,
		speed,
		calories,
	)

	return result, nil
}
