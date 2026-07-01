package domain

import (
	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

// ============================================================================
// CUID HOOKS FOR AUTO-GENERATION
// ============================================================================

// BeforeCreate hook generates a CUID if ID is empty
// Attach this to any model using CUID as primary key

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = cuid.New()
	}
	return nil
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = cuid.New()
	}
	return nil
}

func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = cuid.New()
	}
	return nil
}

func (s *Section) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = cuid.New()
	}
	return nil
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) error {
	if l.ID == "" {
		l.ID = cuid.New()
	}
	return nil
}

func (ld *LessonDetail) BeforeCreate(tx *gorm.DB) error {
	if ld.ID == "" {
		ld.ID = cuid.New()
	}
	return nil
}

func (lp *LessonProgress) BeforeCreate(tx *gorm.DB) error {
	if lp.ID == "" {
		lp.ID = cuid.New()
	}
	return nil
}

func (cert *Certificate) BeforeCreate(tx *gorm.DB) error {
	if cert.ID == "" {
		cert.ID = cuid.New()
	}
	return nil
}

func (d *Discussion) BeforeCreate(tx *gorm.DB) error {
	if d.ID == "" {
		d.ID = cuid.New()
	}
	return nil
}

func (r *Reply) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = cuid.New()
	}
	return nil
}

func (s *Showcase) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = cuid.New()
	}
	return nil
}
func (sc *ShowcaseComment) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == "" {
		sc.ID = cuid.New()
	}
	return nil
}

func (b *Blog) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = cuid.New()
	}
	return nil
}

func (cp *CommunityPoint) BeforeCreate(tx *gorm.DB) error {
	if cp.ID == "" {
		cp.ID = cuid.New()
	}
	return nil
}

func (ul *UserLevel) BeforeCreate(tx *gorm.DB) error {
	if ul.ID == "" {
		ul.ID = cuid.New()
	}
	return nil
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) error {
	if rt.ID == "" {
		rt.ID = cuid.New()
	}
	return nil
}

func (sc *StudyCase) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == "" {
		sc.ID = cuid.New()
	}
	return nil
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = cuid.New()
	}
	return nil
}

func (a *AdminAuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = cuid.New()
	}
	return nil
}

func (f *FAQ) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = cuid.New()
	}
	return nil
}

func (ls *LearningStreak) BeforeCreate(tx *gorm.DB) error {
	if ls.ID == "" {
		ls.ID = cuid.New()
	}
	return nil
}
