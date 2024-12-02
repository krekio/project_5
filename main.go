package main

//Импортнул пакеты
import (
	"fmt"
	"math"
	"time"
)

// Константы для перевода единиц изм.
const (
	MInKm      = 1000
	MinInHours = 60
	LenStep    = 0.65
	CmInM      = 100
)

// Training - структурка для описания тренировки.
type Training struct {
	TrainingType string
	Action       int
	LenStep      float64
	Duration     time.Duration
	Weight       float64
}

// Метод distance(): считаем дистанцию в км.
func (t Training) distance() float64 {
	dist := (float64(t.Action) * t.LenStep) / MInKm
	return dist
}

// Метод meanSpeed(): считаем ср.скорость в км/ч.
func (t Training) meanSpeed() float64 {
	if t.Duration.Hours() == 0 {
		fmt.Println("На ноль делить нельзя!")
	}
	speed := t.distance() / t.Duration.Hours()
	return speed
}

// Calories пока не трогаем, калории для каждой трени считаются по-разному.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage - структурка, в которую запишем полученные Distance,Speed,Calories
type InfoMessage struct {
	TrainingType string
	Duration     time.Duration
	Distance     float64
	Speed        float64
	Calories     float64
}

// TrainingInfo - заполняем структурку InfoMessage полученными значениями.
func (t Training) TrainingInfo() InfoMessage {
	trainingInfo := InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
	return trainingInfo
}

// InfoMessage - метод для вывода инфы в консоль.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator - интерфейс для подсчета калорий для разных типов трень.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

const (
	CaloriesMeanSpeedMultiplier = 18
	CaloriesMeanSpeedShift      = 1.79
)

// Running - встраиваем только Training. Больше ничего не нужно.
type Running struct {
	Training
}

// Calories - метод для подсчета калорий при беге.
func (r Running) Calories() float64 {
	spentCaloriesRunning := (CaloriesMeanSpeedMultiplier*r.meanSpeed() + CaloriesMeanSpeedShift) * r.Weight / MInKm * r.Duration.Hours() * MinInHours
	return spentCaloriesRunning
}

// TrainingInfo - вернули инфу с тренировкой про бег.
func (r Running) TrainingInfo() InfoMessage {
	return r.Training.TrainingInfo()
}

const (
	CaloriesWeightMultiplier      = 0.035
	CaloriesSpeedHeightMultiplier = 0.029
	KmHInMsec                     = 0.278
)

// Walking - добавляем Height для расчетов.
type Walking struct {
	Training
	Height float64
}

// Calories - метод для подсчета калорий при ходьбе.
func (w Walking) Calories() float64 {
	if w.Height == 0 {
		fmt.Println("На ноль делить нельзя!")
	}
	squareSpeed := math.Pow(w.meanSpeed()*KmHInMsec, 2)
	heightInM := w.Height / CmInM
	spentCaloriesWalking := (CaloriesWeightMultiplier*w.Weight + (squareSpeed/heightInM)*CaloriesSpeedHeightMultiplier*w.Weight) * w.Duration.Hours() * MinInHours
	return spentCaloriesWalking
}

// TrainingInfo - вернули инфу с тренировкой про ходьбу.
func (w Walking) TrainingInfo() InfoMessage {
	return w.Training.TrainingInfo()
}

const (
	SwimmingLenStep                  = 1.38
	SwimmingCaloriesMeanSpeedShift   = 1.1
	SwimmingCaloriesWeightMultiplier = 2
)

// Swimming - добавляем LengthPool и CountPool для расчетов.
type Swimming struct {
	Training
	LengthPool int
	CountPool  int
}

// distance - переопределил новую дистанцию для плавания.
func (s Swimming) distance() float64 {
	distance := float64(s.LengthPool*s.CountPool) / MInKm
	return distance
}

// meanSpeed - считаем скорость во время плавания.
func (s Swimming) meanSpeed() float64 {
	if s.Duration.Hours() == 0 {
		fmt.Println("На ноль делить нельзя!")
	}
	swimmingSpeed := float64(s.LengthPool*s.CountPool) / MInKm / s.Duration.Hours()
	return swimmingSpeed
}

// TrainingInfoForSwimming - TrainingInfo, только для плавания
func (s Swimming) TrainingInfoForSwimming() InfoMessage {
	trainingInfo := InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     s.distance(),
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
	return trainingInfo
}

// Calories - метод для подсчета калорий при плавании.
func (s Swimming) Calories() float64 {
	spentCaloriesSwimming := (s.meanSpeed() + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
	return spentCaloriesSwimming
}

// TrainingInfo - вернули инфу с тренировкой про плавание.
func (s Swimming) TrainingInfo() InfoMessage {
	return s.TrainingInfoForSwimming()
}

// ReadData - работает для всех типов тренировок. Получаем кол-во калорий. В инфо кладем структурку о тренировке и присваиваем info.Calories посчитанные калории.
func ReadData(training CaloriesCalculator) string {
	calories := training.Calories()
	info := training.TrainingInfo()
	info.Calories = calories
	return fmt.Sprint(info)
}
func main() {
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}
	fmt.Println(ReadData(swimming))
	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}
	fmt.Println(ReadData(walking))
	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}
	fmt.Println(ReadData(running))
}
